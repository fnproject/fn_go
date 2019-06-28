package oracle

import (
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"path"

	"github.com/fnproject/fn_go/client/version"
	"github.com/fnproject/fn_go/clientv2"
	"github.com/fnproject/fn_go/provider"
	openapi "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/oracle/oci-go-sdk/common"
	oci "github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/common/auth"
)

const (
	CfgPassPhrase       = "oracle.pass-phrase"
	CfgCompartmentID    = "oracle.compartment-id"
	CfgDisableCerts     = "oracle.disable-certs"
	CompartmentMetadata = "http://169.254.169.254/opc/v1/instance/compartmentId"
)

// Provider :  Oracle Instance Principal Authentication provider
// This is used to authenticate using the Instance Principal token mechanism available from OCI host instances.
type Provider struct {
	// FnApiUrl is the endpoint to use for API interactions
	FnApiUrl *url.URL
	// FnCallUrl is the endpoint used for call interactions
	FnCallUrl *url.URL
	// IPProvider is the instance principal ConfigurationProvider
	IPProvider common.ConfigurationProvider
	//DisableCerts indicates if server certificates should be ignored
	DisableCerts bool
	//CompartmentID is the ocid of the functions compartment ID for a given function
	CompartmentID string
}

func (op *Provider) APIClientv2() *clientv2.Fn {
	runtime := openapi.New(op.FnApiUrl.Host, path.Join(op.FnApiUrl.Path, clientv2.DefaultBasePath), []string{op.FnApiUrl.Scheme})
	runtime.Transport = op.WrapCallTransport(runtime.Transport)
	return clientv2.New(runtime, strfmt.Default)
}

type Response struct {
	Annotations Annotations `json:"annotations"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	Name        string      `json:"name"`
}

type Annotations struct {
	CompartmentID string `json:"oracle.com/oci/compartmentId"`
	ShortCode     string `json:"oracle.com/oci/appCode"`
}

func NewFromConfig(configSource provider.ConfigSource, passphraseSource provider.PassPhraseSource) (provider.Provider, error) {
	ip, err := auth.InstancePrincipalConfigurationProvider()
	if err != nil {
		return nil, err
	}

	cfgApiUrl := configSource.GetString(provider.CfgFnAPIURL)
	if cfgApiUrl == "" {
		region, err := ip.Region()
		if err != nil {
			return nil, err
		}
		// Construct the API endpoint from the "nearby" endpoint
		cfgApiUrl = fmt.Sprintf("https://functions.%s.oraclecloud.com", region)
	}
	apiUrl, err := provider.CanonicalFnAPIUrl(cfgApiUrl)
	if err != nil {
		return nil, err
	}

	compartmentID := configSource.GetString(CfgCompartmentID)
	if compartmentID == "" {
		// Get the local compartment ID from the metadata endpoint
		resp, err := http.DefaultClient.Get(CompartmentMetadata)
		if err != nil {
			return nil, fmt.Errorf("problem fetching compartment Id from metadata endpoint %s", err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("problem fetching compartment Id from metadata endpoint %s", err)
		}
		compartmentID = string(body)
	}
	return &Provider{FnApiUrl: apiUrl,
		IPProvider:    ip,
		DisableCerts:  configSource.GetBool(CfgDisableCerts),
		CompartmentID: compartmentID,
	}, nil
}

func (op *Provider) APIURL() *url.URL {
	return op.FnApiUrl
}

func (op *Provider) WrapCallTransport(roundTripper http.RoundTripper) http.RoundTripper {
	if op.DisableCerts {
		roundTripper = InsecureRoundTripper(roundTripper)
	}

	ociClient := common.RequestSigner(op.IPProvider, []string{"host", "date", "(request-target)"}, []string{"content-length", "content-type", "x-content-sha256"})

	signingRoundTrripper := ociSigningRoundTripper{
		transport: roundTripper,
		ociClient: ociClient,
	}

	roundTripper = compartmentIDRoundTripper{
		transport:     signingRoundTrripper,
		compartmentID: op.CompartmentID,
	}

	roundTripper = requestIdRoundTripper{
		transport: roundTripper,
	}

	return roundTripper
}

func (op *Provider) APIClient() *clientv2.Fn {
	runtime := openapi.New(op.FnApiUrl.Host, path.Join(op.FnApiUrl.Path, clientv2.DefaultBasePath), []string{op.FnApiUrl.Scheme})
	runtime.Transport = op.WrapCallTransport(runtime.Transport)
	return clientv2.New(runtime, strfmt.Default)
}

func (op *Provider) VersionClient() *version.Client {
	runtime := openapi.New(op.FnApiUrl.Host, op.FnApiUrl.Path, []string{op.FnApiUrl.Scheme})
	runtime.Transport = op.WrapCallTransport(runtime.Transport)
	return version.New(runtime, strfmt.Default)
}

type ociKeyProvider struct {
	ID  string
	key *rsa.PrivateKey
}

func (kp ociKeyProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	return kp.key, nil
}

func (kp ociKeyProvider) KeyID() (string, error) {
	return kp.ID, nil
}

type ociSigningRoundTripper struct {
	ociClient common.HTTPRequestSigner
	transport http.RoundTripper
}

func (t ociSigningRoundTripper) RoundTrip(request *http.Request) (response *http.Response, err error) {
	if request.Header.Get("Date") == "" {
		request.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}
	request, err = signRequest(t.ociClient, request)

	if err != nil {
		return
	}

	response, err = t.transport.RoundTrip(request)

	return
}

// Add the necessary headers and sign the request
func signRequest(signer common.HTTPRequestSigner, request *http.Request) (signedRequest *http.Request, err error) {
	// Check that a Date header is set, otherwise authentication will fail
	if request.Header.Get("Date") == "" {
		return nil, fmt.Errorf("Date header must be present and non-empty on request")
	}
	if request.Method == "POST" || request.Method == "PATCH" || request.Method == "PUT" {
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Length", fmt.Sprintf("%d", request.ContentLength))
	}

	err = signer.Sign(request)

	return request, err
}

// http.RoundTripper middleware that injects an opc-request-id header
type requestIdRoundTripper struct {
	transport http.RoundTripper
}

func (t requestIdRoundTripper) RoundTrip(request *http.Request) (response *http.Response, e error) {
	requestID := provider.GetRequestID(request.Context())
	if requestID != "" {
		request.Header.Set("Opc-Request-Id", requestID)

	}
	response, e = t.transport.RoundTrip(request)
	return
}

//  http.RoundTripper middleware that adds an opc-compartment-id header to all requests
type compartmentIDRoundTripper struct {
	transport     http.RoundTripper
	compartmentID string
}

func (t compartmentIDRoundTripper) RoundTrip(request *http.Request) (response *http.Response, e error) {
	request.Header.Set("opc-compartment-id", t.compartmentID)
	response, e = t.transport.RoundTrip(request)
	return
}

// Skip verification of insecure certs
func InsecureRoundTripper(roundTripper http.RoundTripper) http.RoundTripper {
	transport := roundTripper.(*http.Transport)
	if transport != nil {
		if transport.TLSClientConfig != nil {
			transport.TLSClientConfig.InsecureSkipVerify = true
		} else {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}

	return transport
}

func privateKey(config provider.ConfigSource, passphrase provider.PassPhraseSource, pkeyFilePath string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(pkeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load private key from file: %s. Error: %s \n", pkeyFilePath, err)
	}

	var key *rsa.PrivateKey
	pKeyPword := config.GetString(CfgPassPhrase)
	if !config.IsSet(CfgPassPhrase) {
		if key, err = getPrivateKey(keyBytes, pKeyPword, pkeyFilePath); key == nil {
			pKeyPword, err = passphrase.ChallengeForPassPhrase("oracle.privateKey", fmt.Sprintf("Enter passphrase for private key %s", pkeyFilePath))
			if err != nil {
				return nil, err
			}
		}
	}
	key, err = getPrivateKey(keyBytes, pKeyPword, pkeyFilePath)
	return key, err
}

func getPrivateKey(keyBytes []byte, pKeyPword, pkeyFilePath string) (*rsa.PrivateKey, error) {
	key, err := oci.PrivateKeyFromBytes(keyBytes, oci.String(pKeyPword))
	if err != nil {
		if pKeyPword != "" {
			return nil, fmt.Errorf("Unable to load private key from file bytes: %s. Error: %s \n", pkeyFilePath, err)
		}
	}

	return key, nil
}
