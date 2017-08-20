package stripe

import (
	"fmt"
	"strings"

	stripe "github.com/stripe/stripe-go"

	"github.com/BTBurke/recur/pb"
	context "golang.org/x/net/context"
)

type StripePlanClient struct {
	Key string
}

func (p *StripePlanClient) Create(ctx context.Context, req *pb.CreatePlanRequest) (*pb.PlanResponse, error) {
	return nil, nil
}

func marshalPlanResponse(plan *stripe.Plan) *pb.PlanResponse {
	return &pb.PlanResponse{
		Responses: &pb.PlanResponse_Response{
			Response: &pb.PlanSuccess{
				Id:       plan.ID,
				Object:   "plan",
				Livemode: plan.Live,
				Amount:   plan.Amount,
				Created:  plan.Created,
				Currency: getCurrencyEnum(plan.Currency),
				Interval: getIntervalEnum(plan.Interval),
			},
		},
	}
}

// constant conversions from stripe to protobuf - currency
func getCurrencyEnum(c stripe.Currency) pb.Currency {
	return pb.Currency(pb.Currency_value[strings.ToUpper(fmt.Sprintf("%s", c))])
}

// constant conversions from stripe to protobuf - plan interval
func getIntervalEnum(i stripe.PlanInterval) pb.Interval {
	return pb.Interval(pb.Interval_value[strings.ToLower(fmt.Sprintf("%s", i))])
}
