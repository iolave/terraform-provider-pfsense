package pfsenseclient

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/hashicorp/go-cleanhttp"
)

func NewHTTPClient(TLSSkipVerify bool) *http.Client {
	jar, err := cookiejar.New(nil)

	if err != nil {
		panic(err)
	}

	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: TLSSkipVerify,
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}

	return client
}

// Sends an http request using the net/http package.
// In case of a request failure, n-request retries
// will be sent in order to retrieve a response.
//
// When maxAttempts is a negative number, maxAttempts
// will be assumed as zero.
func DoWithRetries(client *http.Client, req *http.Request, maxAttempts int) (*http.Response, error) {
	count := 0

	// Sets maxAttempts as zero when a negative number
	// is given
	if maxAttempts < 0 {
		maxAttempts = 0
	}

	// Iterates maxAttempts times plus the one to include
	// the first call before retrying in case of failure
	for attempt := 0; attempt < maxAttempts+1; attempt++ {
		count++
		res, err := client.Do(req)

		// When the request's context contains an error we
		// won't retry cuz the request was cancelled at some
		// point.
		if ctxErr := req.Context().Err(); ctxErr != nil {
			return nil, ctxErr
		}

		// In the http.Do method an error is returned 'by client
		// policy (such as CheckRedirect), or failure to speak HTTP
		// (such as a network connectivity problem)' which means
		// that we'll ignore the error.
		if err != nil {
			continue
		}

		// At this point we need to tell wether to retry or not
		// based on the response content.
		if res.StatusCode == 0 {
			continue
		}

		if res.StatusCode >= 500 && res.StatusCode != http.StatusNotImplemented {
			continue
		}

		return res, nil
	}

	return nil, errors.New("neither error nor response found")
}
