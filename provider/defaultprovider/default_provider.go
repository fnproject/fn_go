package defaultprovider

import (
	openapi "github.com/go-openapi/runtime/client"

	"fmt"
	"github.com/fnproject/fn_go/client"
	"github.com/fnproject/fn_go/client/version"
	"github.com/fnproject/fn_go/provider"
	"github.com/go-openapi/strfmt"
	"net/http"
	"net/url"
)

const CfgFnToken = "token"
const CfgFnUser = "user"
const CfgFnPasswd = "password"


// Provider is the default Auth provider
type Provider struct {
	// Optional token to add as  bearer token to auth calls
	Token string
	// Optional basic auth user
	User string
	// Optional basic auth password
	Password string
	// API url to use for FN API interactions
	FnApiUrl *url.URL
	// URL to use for FN call interactions
	CallUrl *url.URL
}

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
		Token:    configSource.GetString(CfgFnToken),
		User:    configSource.GetString(CfgFnUser),
		Password:    configSource.GetString(CfgFnPasswd),
		FnApiUrl: apiUrl,
		CallUrl:  callUrl,
	}, nil
}

func (dp *Provider) WrapCallTransport(t http.RoundTripper) http.RoundTripper {
	return t
}

func (dp *Provider) CallURL() *url.URL {
	return dp.CallUrl
}
func (dp *Provider) APIURL() *url.URL {
	return dp.FnApiUrl
}

func (dp *Provider) APIClient() *client.Fn {
	transport := openapi.New(dp.FnApiUrl.Host, dp.FnApiUrl.Path, []string{dp.FnApiUrl.Scheme})
	if dp.Token != "" {
		transport.DefaultAuthentication = openapi.BearerToken(dp.Token)
	} else if dp.User != "" {
		transport.DefaultAuthentication = openapi.BasicAuth(dp.User, dp.Password)
	}

	return client.New(transport, strfmt.Default)
}

func (op *Provider) VersionClient() *version.Client {
	runtime := openapi.New(op.FnApiUrl.Host, "/", []string{op.FnApiUrl.Scheme})
	runtime.Transport = op.WrapCallTransport(runtime.Transport)
	return version.New(runtime, strfmt.Default)
}
