package stripe

import (
	"github.com/stripe/stripe-go"
	context "golang.org/x/net/context"
)

func paramsFromContext(ctx context.Context, key string, meta *map[string]string) stripe.Params {
	p := stripe.Params{}
	if meta != nil {
		for k, v := range *meta {
			p.AddMeta(k, v)
		}
	}
	headers, ok := ctx.Value("headers").(map[string]string)
	if ok {
		for k, v := range headers {
			p.Headers.Add(k, v)
		}
	}
	idempotencyKey, ok := ctx.Value("idempotency").(string)
	if ok {
		p.IdempotencyKey = idempotencyKey
	}
	stripeAccount, ok := ctx.Value("connectkey").(string)
	if ok {
		p.SetStripeAccount(stripeAccount)
	}

	return p
}
