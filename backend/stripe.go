package backend

import (
	"github.com/BTBurke/recur/backend/stripe"
	log "github.com/sirupsen/logrus"
)

func NewStripePlanClient(key string, logger log.StdLogger) *stripe.StripePlanClient {
	return &stripe.StripePlanClient{
		Key:    key,
		Logger: logger,
	}
}
