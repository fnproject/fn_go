package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/fnproject/fn_go/client/version"
	"github.com/fnproject/fn_go/clientv2"
)

const MaximumRequestBodySize = 5 * 1024 * 1024

// ProviderFunc constructs a provider
type ProviderFunc func(config ConfigSource, source PassPhraseSource) (Provider, error)

//Providers describes a set of providers
type Providers struct {
	Providers map[string]ProviderFunc
}

// Register adds a named provider to a configuration
func (c *Providers) Register(name string, pf ProviderFunc) {
	if len(c.Providers) == 0 {
		c.Providers = make(map[string]ProviderFunc)
	}
	c.Providers[name] = pf
}

// Provider creates API clients for Fn calls adding any required middleware
type Provider interface {
	// APIURL returns the current API URL base to use with this provider
	APIURL() *url.URL
	// WrapCallTransport adds any request signing or auth to an existing round tripper for calls
	WrapCallTransport(http.RoundTripper) http.RoundTripper
	APIClientv2() *clientv2.Fn
	VersionClient() *version.Client
	Invoke(invokeURL string, content io.Reader, output io.Writer, headers http.Header, contentType string, debug bool) error
}

// CanonicalFnAPIUrl canonicalises an *FN_API_URL  to a default value
func CanonicalFnAPIUrl(urlStr string) (*url.URL, error) {
	if !strings.Contains(urlStr, "://") {
		urlStr = fmt.Sprint("http://", urlStr)
	}

	parseUrl, err := url.Parse(urlStr)

	if err != nil {
		return nil, fmt.Errorf("unparsable FN API Url: %s. Error: %s", urlStr, err)
	}

	if parseUrl.Port() == "" {
		if parseUrl.Scheme == "http" {
			parseUrl.Host = fmt.Sprint(parseUrl.Host, ":80")
		}
		if parseUrl.Scheme == "https" {
			parseUrl.Host = fmt.Sprint(parseUrl.Host, ":443")
		}
	}

	//Remove /v1 from any paths here internal URL is now base URL

	if strings.HasSuffix(parseUrl.Path, "/v1") {
		parseUrl.Path = strings.TrimSuffix(parseUrl.Path, "v1")
	} else if strings.HasSuffix(parseUrl.Path, "/v1/") {
		parseUrl.Path = strings.TrimSuffix(parseUrl.Path, "v1/")
	}

	return parseUrl, nil
}

//ProviderFromConfig returns the provider corresponding to a given identifier populated with configuration from source - if a passphrase is required then it is request from phraseSource
func (c *Providers) ProviderFromConfig(id string, source ConfigSource, phraseSource PassPhraseSource) (Provider, error) {
	p, ok := c.Providers[id]
	if !ok {
		return nil, fmt.Errorf("No provider with id  '%s' is registered", id)
	}
	return p(source, phraseSource)
}

type apiErr struct {
	Message string `json:"message"`
}
type callID struct {
	CallID string `json:"call_id"`
	Error  apiErr `json:"error"`
}

func Invoke(provider Provider, invokeURL string, content io.Reader, output io.Writer, headers http.Header, contentType string, debug bool) error {

	method := "POST"

	// Read the request body (up to the maximum size), as this is used in the
	// authentication signature
	var req *http.Request
	if content != nil {
		b, err := ioutil.ReadAll(io.LimitReader(content, MaximumRequestBodySize))
		if err != nil {
			return err
		}
		buffer := bytes.NewBuffer(b)
		req, err = http.NewRequest(method, invokeURL, buffer)
		if err != nil {
			return fmt.Errorf("Error creating request to service: %s", err)
		}
	} else {
		var err error
		req, err = http.NewRequest(method, invokeURL, nil)
		if err != nil {
			return fmt.Errorf("Error creating request to service: %s", err)
		}
	}

	for key, v := range headers {
		for _, value := range v {
			req.Header.Add(key, value)
		}
	}

	transport := provider.WrapCallTransport(http.DefaultTransport)
	httpClient := http.Client{Transport: transport}

	if debug {
		b, err := httputil.DumpRequestOut(req, content != nil)
		if err != nil {
			return err
		}
		fmt.Printf(string(b) + "\n")
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("Error invoking fn: %s", err)
	}

	if debug {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		fmt.Printf(string(b) + "\n")
	}

	// for sync calls
	if call_id, found := resp.Header["FN_CALL_ID"]; found {
		if debug {
			fmt.Fprint(os.Stderr, fmt.Sprintf("Call ID: %v\n", call_id[0]))
		}
		io.Copy(output, resp.Body)
	} else {
		// for async calls and error discovering
		c := &callID{}
		err = json.NewDecoder(resp.Body).Decode(c)
		if err == nil {
			// decode would not fail in both cases:
			// - call id in body
			// - error in body
			// that's why we need to check values of attributes
			if c.CallID != "" {
				fmt.Fprint(os.Stderr, fmt.Sprintf("Call ID: %v\n", c.CallID))
			} else {
				fmt.Fprint(output, fmt.Sprintf("Error: %v\n", c.Error.Message))
			}
		} else {
			return err
		}
	}

	if resp.StatusCode >= 400 {
		// TODO: parse out error message
		return fmt.Errorf("Error calling function: status %v", resp.StatusCode)
	}

	return nil
}
