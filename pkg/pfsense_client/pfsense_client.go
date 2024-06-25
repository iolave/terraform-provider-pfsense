package pfsenseclient

import (
	"fmt"
	"net/url"
	"time"

	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

const (
	DefaultURL           = "https://192.168.1.1"
	DefaultUsername      = "admin"
	DefaultTLSSkipVerify = false
	DefaultRetryMinWait  = time.Second
	DefaultRetryMaxWait  = 5 * time.Second
	DefaultMaxAttempts   = 3
)

func New(options Options) (*PfsenseClient, error) {

	// If option's url is not defined (empty)
	// it will set the default pfsense url.
	if options.URL.String() == "" {
		url, err := url.Parse(DefaultURL)

		if err != nil {
			return nil, err
		}

		options.URL = url
	}

	// If option's username is not defined (empty)
	// it will set the default pfsense username.
	if options.Username == "" {
		options.Username = DefaultUsername
	}

	// When no password is given an error will be
	// returned
	if options.Password == "" {
		return nil, fmt.Errorf("%w, password required", pfsense.ErrClientValidation)
	}

	if options.TLSSkipVerify == nil {
		b := DefaultTLSSkipVerify
		options.TLSSkipVerify = &b
	}

	if options.RetryMinWait == nil {
		td := DefaultRetryMinWait
		options.RetryMinWait = &td
	}

	if options.RetryMaxWait == nil {
		td := DefaultRetryMaxWait
		options.RetryMaxWait = &td
	}

	if options.MaxAttempts == nil {
		i := DefaultMaxAttempts
		options.MaxAttempts = &i
	}

	client := &PfsenseClient{
		Options: &options,
	}

	return client, nil
}

type PfsenseClient struct {
	Options *Options
}

type Options struct {
	URL           *url.URL
	Username      string
	Password      string
	TLSSkipVerify *bool
	RetryMinWait  *time.Duration
	RetryMaxWait  *time.Duration
	MaxAttempts   *int
}
