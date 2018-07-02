package defaultprovider

import (
	openapi "github.com/go-openapi/runtime/client"

	"fmt"
	"net/http"
	"net/url"

	"github.com/fnproject/fn_go/client"
	"github.com/fnproject/fn_go/client/version"
	"github.com/fnproject/fn_go/provider"
	"github.com/go-openapi/strfmt"
	"github.com/fnproject/fn_go/clientv2"
	"path"
	)

// Provider is the default Auth provider
type Provider struct {
	// Optional token to add as  bearer token to auth calls
	Token string
	// API url to use for FN API interactions
	FnApiUrl *url.URL
	// URL to use for FN call interactions
	CallUrl *url.URL
}

func (dp *Provider) APIClientv2() *clientv2.Fn {
	transport := openapi.New(dp.FnApiUrl.Host, path.Join(dp.FnApiUrl.Path ,clientv2.DefaultBasePath), []string{dp.FnApiUrl.Scheme})
	if dp.Token != "" {
		transport.DefaultAuthentication = openapi.BearerToken(dp.Token)
	}

	return clientv2.New(transport, strfmt.Default)}

//  NewFromConfig creates a default provider  that does un-authenticated calls to
func NewFromConfig(configSource provider.ConfigSource, _ provider.PassPhraseSource) (provider.Provider, error) {

	apiUrl, err := provider.CanonicalFnAPIUrl(configSource.GetString(provider.CfgFnAPIURL))

	if err != nil {
		return nil, err
	}

	var callUrl *url.URL
	callUrlStr := configSource.GetString(provider.CfgFnCallURL)
	if callUrlStr == "" {
		myCallUrl := *apiUrl
		callUrl = &myCallUrl
		callUrl.Path = "/r"
	} else {
		callUrl, err = url.Parse(callUrlStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse call url %s specified in %s: %s", callUrlStr, provider.CfgFnCallURL, err)
		}
	}
	return &Provider{
		Token:    configSource.GetString(provider.CfgFnToken),
		FnApiUrl: apiUrl,
		CallUrl:  callUrl,
	}, nil
}

func (dp *Provider) WrapCallTransport(t http.RoundTripper) http.RoundTripper {
	return t
}

func (dp *Provider) CallURL(appName string) (*url.URL, error) {
	return dp.CallUrl, nil
}
func (dp *Provider) APIURL() *url.URL {
	return dp.FnApiUrl
}

func (dp *Provider) APIClient() *client.Fn {
	join := path.Join(dp.FnApiUrl.Path, client.DefaultBasePath)
	transport := openapi.New(dp.FnApiUrl.Host, join, []string{dp.FnApiUrl.Scheme})
	if dp.Token != "" {
		transport.DefaultAuthentication = openapi.BearerToken(dp.Token)
	}

	return client.New(transport, strfmt.Default)
}

func (op *Provider) VersionClient() *version.Client {
	runtime := openapi.New(op.FnApiUrl.Host, "/", []string{op.FnApiUrl.Scheme})
	runtime.Transport = op.WrapCallTransport(runtime.Transport)
	return version.New(runtime, strfmt.Default)
}
