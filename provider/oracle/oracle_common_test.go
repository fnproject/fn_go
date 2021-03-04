package oracle

import (
	"net/http"
	"testing"
)

func TestInsecureRoundTripper(t *testing.T) {
	testCases := []http.RoundTripper{
		http.DefaultTransport, // Normal case
		nil,                   // Ensure it creates a DefaultTransport in the nil case
	}

	for _, transport := range testCases {
		roundTripper := InsecureRoundTripper(transport)
		transport, ok := roundTripper.(*http.Transport)

		if !ok {
			t.Fatal("Transport not correctly returned")
		}

		if transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
			t.Fatal("InsecureSkipVerify not correctly set on transport")
		}
	}
}
