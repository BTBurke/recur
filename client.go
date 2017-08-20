package recur

import (
	"github.com/BTBurke/recur/backend"
)

// ClientType enumerates possible backends services (Stripe only for now)
type ClientType int

const (
	// Use Stripe as the backend for recurring billing
	StripeClient ClientType = iota
)

type RunMode int

const (
	RunAsService RunMode = iota
	RunAsLibrary
)

// Client
type Client struct {
	Key     string
	Backend ClientType
	RunMode RunMode
	Plan    *PlanClient
}

// ClientOption is a function that applies an option to the client configuration
type ClientOption func(c *Client) error

// NewClient returns a new client using the chosen service (e.g. Stripe) as the backend. Call this
// to create a client when using recur as a library in your own Go project.  Use ClientOption
// to configure optional behavior such as logging (disabled by default).
func NewClient(service ClientType, key string, opts ...ClientOption) *Client {
	return newClient(service, key, RunAsLibrary)
}

// NewGRPCClient returns a new client using the chosen service (e.g. Stripe) as the backend. Call this
// to create a client when using recur as a microservice.  Use ClientOption
// to configure optional behavior.  Logging is enabled at INFO by default when running as a service.
func NewGRPCClient(service ClientType, key string) *Client {
	return newClient(service, key, RunAsService)
}

func newClient(service ClientType, key string, run RunMode) *Client {
	switch service {
	case StripeClient:
		return &Client{
			Key:     key,
			RunMode: run,
			Plan:    &PlanClient{backend: backend.NewStripePlanClient(key)},
		}
	default:
		return nil
	}
}

// func NoLog() ClientOption {
// 	return func(c *Client) error {
// 		c.Logger =
// 	}
// }
