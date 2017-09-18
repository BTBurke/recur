package recur

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/BTBurke/recur/backend/stripe"
	log "github.com/sirupsen/logrus"
)

// ClientType enumerates possible backends services (Stripe only for now)
type ClientType int

const (
	// Use Stripe as the backend for recurring billing
	StripeClient ClientType = iota
)

type runMode int

const (
	runAsService runMode = iota
	runAsLibrary
)

// Client
type Client struct {
	Key     string
	Backend ClientType
	Timeout time.Duration
	Logger  *log.Logger

	Plan *PlanClient

	runMode runMode
}

// ClientOption is a function that applies an option to the client configuration
type ClientOption func(c *Client) error

// NewClient returns a new client using the chosen service (e.g. Stripe) as the backend. Call this
// to create a client when using recur as a library in your own Go project.  Use ClientOption
// to configure optional behavior such as logging (disabled by default).
func NewClient(service ClientType, key string, opts ...ClientOption) (*Client, error) {
	return newClient(service, key, runAsLibrary, opts...)
}

// NewGRPCClient returns a new client using the chosen service (e.g. Stripe) as the backend. Call this
// to create a client when using recur as a microservice.  Use ClientOption
// to configure optional behavior.  Logging is enabled at INFO by default to Stdout.
func NewGRPCClient(service ClientType, key string, opts ...ClientOption) (*Client, error) {
	return newClient(service, key, runAsService, opts...)
}

func newClient(service ClientType, key string, run runMode, opts ...ClientOption) (*Client, error) {

	c := &Client{
		Key:     key,
		runMode: run,
		Logger:  log.New(),
	}

	defaultOpts := []ClientOption{
		LogLevel(LogLevelInfo),
		LogFormat(TextFormatter),
	}
	switch run {
	case runAsLibrary:
		defaultOpts = append(defaultOpts, NoLog())
	default:
		defaultOpts = append(defaultOpts, LogOutput(os.Stdout))
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	switch service {
	case StripeClient:
		c.Plan = &PlanClient{backend: stripe.NewPlanClient(key, c.Logger)}
		return c, nil
	default:
		return nil, fmt.Errorf("unknown backend service")
	}
}

func NoLog() ClientOption {
	return func(c *Client) error {
		c.Logger.Out = ioutil.Discard
		return nil
	}
}

type LoggerLevel int

const (
	LogLevelDebug LoggerLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func LogLevel(level LoggerLevel) ClientOption {
	return func(c *Client) error {
		switch level {
		case LogLevelDebug:
			c.Logger.Level = log.DebugLevel
		case LogLevelInfo:
			c.Logger.Level = log.InfoLevel
		case LogLevelWarn:
			c.Logger.Level = log.WarnLevel
		case LogLevelError:
			c.Logger.Level = log.ErrorLevel
		}
		return nil
	}
}

type LoggerFormat int

const (
	JSONFormatter LoggerFormat = iota
	TextFormatter
)

func LogFormat(f LoggerFormat) ClientOption {
	return func(c *Client) error {
		switch f {
		case JSONFormatter:
			c.Logger.Formatter = new(log.JSONFormatter)
		case TextFormatter:
			c.Logger.Formatter = new(log.TextFormatter)
		}
		return nil
	}
}

func LogOutput(w io.Writer) ClientOption {
	return func(c *Client) error {
		c.Logger.Out = w
		return nil
	}
}
