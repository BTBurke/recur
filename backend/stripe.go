package backend

import (
	"github.com/BTBurke/recur/backend/stripe"
)

func NewStripePlanClient(key string) *stripe.StripePlanClient {
	return &stripe.StripePlanClient{Key: key}
}
