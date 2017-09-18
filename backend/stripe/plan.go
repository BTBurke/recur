package stripe

import (
	"github.com/BTBurke/recur/backend"
	"github.com/BTBurke/recur/pb"
	"github.com/cenkalti/backoff"
	"github.com/stripe/stripe-go/plan"

	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	context "golang.org/x/net/context"
)

// interface for the Stripe plan API
type planClient interface {
	New(params *stripe.PlanParams) (*stripe.Plan, error)
	Get(id string, params *stripe.PlanParams) (*stripe.Plan, error)
	Update(id string, params *stripe.PlanParams) (*stripe.Plan, error)
	Del(id string, params *stripe.PlanParams) (*stripe.Plan, error)
	List(params *stripe.PlanListParams) *plan.Iter
}

type StripePlanClient struct {
	key    string
	logger log.StdLogger

	// api allows mocking the Stripe backend
	api planClient
}

func NewPlanClient(key string, logger log.StdLogger) *StripePlanClient {
	return &StripePlanClient{
		key:    key,
		logger: logger,
		api: plan.Client{
			B:   stripe.GetBackend(stripe.SupportedBackend("api")),
			Key: key,
		},
	}
}

func (p *StripePlanClient) Create(ctx context.Context, req *pb.CreatePlanRequest) (*pb.PlanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	planParams := planCreateToPlanParams(ctx, p.key, req)

	resp := new(pb.PlanResponse)
	err := backoff.Retry(
		retryablePlan(planParams, p.api, resp, planCreate),
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
	)

	return resp, err
}

func (p *StripePlanClient) Update(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.PlanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	params := planUpdateToPlanParams(ctx, p.key, req)

	resp := new(pb.PlanResponse)
	err := backoff.Retry(
		retryablePlan(params, p.api, resp, planUpdate),
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
	)

	return resp, err
}

func (p *StripePlanClient) Delete(ctx context.Context, req *pb.DeletePlanRequest) (*pb.DeletePlanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	params := planDeleteToPlanParams(ctx, p.key, req)

	resp := new(pb.DeletePlanResponse)
	err := backoff.Retry(
		retryablePlanDelete(params, p.api, resp, planDelete),
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
	)

	return resp, err
}

func (p *StripePlanClient) Get(ctx context.Context, req *pb.GetPlanRequest) (*pb.PlanResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	params := planGetToPlanParams(ctx, p.key, req)

	resp := new(pb.PlanResponse)
	err := backoff.Retry(
		retryablePlan(params, p.api, resp, planGet),
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
	)

	return resp, err
}

// planStreamer implements the PlanStreamer interface, converting Stripe responses
// to a PlanResponse.
type planStreamer struct {
	iter *plan.Iter
}

func (s *planStreamer) Next() bool {
	return s.iter.Next()
}

func (s *planStreamer) Current() *pb.PlanResponse {
	switch {
	case s.iter.Err() != nil:
		return respToPlanError(s.iter.Err().(*stripe.Error))
	default:
		return respToPlanSuccess(s.iter.Plan())
	}
}

func (p *StripePlanClient) List(ctx context.Context, req *pb.ListPlansRequest) (backend.PlanStreamer, error) {

	params := planListToListParams(ctx, p.key, req)

	streamer := new(planStreamer)
	err := backoff.Retry(
		retryablePlanList(params, p.api, streamer),
		backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
	)

	return streamer, err
}
