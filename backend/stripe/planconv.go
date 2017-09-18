package stripe

import (
	"fmt"
	"strings"

	"github.com/BTBurke/recur/pb"
	"github.com/stripe/stripe-go"
	context "golang.org/x/net/context"
)

// convert from a plan create request to PlanParams
func planCreateToPlanParams(ctx context.Context, key string, req *pb.CreatePlanRequest) *stripe.PlanParams {
	return &stripe.PlanParams{
		Params:        paramsFromContext(ctx, key, &req.Metadata),
		ID:            req.Id,
		Name:          req.Name,
		Currency:      pbToStripeCurrency(req.GetCurrency()),
		Amount:        req.Amount,
		Interval:      pbToStripeInterval(req.GetInterval()),
		IntervalCount: defaultUint64(req.IntervalCount, 1),
		TrialPeriod:   req.TrialPeriodDays,
		Statement:     req.StatementDescriptor,
	}
}

// convert from plan update to PlanParams
func planUpdateToPlanParams(ctx context.Context, key string, req *pb.UpdatePlanRequest) *stripe.PlanParams {
	return &stripe.PlanParams{
		Params:      paramsFromContext(ctx, key, &req.Metadata),
		ID:          req.Id,
		Name:        req.Name,
		Statement:   req.StatementDescriptor,
		TrialPeriod: req.TrialPeriodDays,
	}
}

// convert from plan delete to PlanParams
func planDeleteToPlanParams(ctx context.Context, key string, req *pb.DeletePlanRequest) *stripe.PlanParams {
	return &stripe.PlanParams{
		Params: paramsFromContext(ctx, key, nil),
		ID:     req.Id,
	}
}

// convert from plan get to PlanParams
func planGetToPlanParams(ctx context.Context, key string, req *pb.GetPlanRequest) *stripe.PlanParams {
	return &stripe.PlanParams{
		Params: paramsFromContext(ctx, key, nil),
		ID:     req.Id,
	}
}

func planListToListParams(ctx context.Context, key string, req *pb.ListPlansRequest) *stripe.PlanListParams {
	switch {
	case req == nil:
		return &stripe.PlanListParams{
			ListParams: stripe.ListParams{
				Limit: 10,
			},
		}
	default:
		return &stripe.PlanListParams{
			ListParams: stripe.ListParams{
				Start: req.StartingAfter,
				End:   req.EndingBefore,
				Limit: defaultInt(int(req.Limit), 10),
			},
			CreatedRange: &stripe.RangeQueryParams{
				GreaterThan:        req.Created.Gt,
				GreaterThanOrEqual: req.Created.Gte,
				LesserThan:         req.Created.Lt,
				LesserThanOrEqual:  req.Created.Lte,
			},
		}
	}
}

// convert a success response from Stripe to a PlanResponse (success)
func respToPlanSuccess(plan *stripe.Plan) *pb.PlanResponse {
	return &pb.PlanResponse{
		Responses: &pb.PlanResponse_Success{
			Success: &pb.Plan{
				Id:                  plan.ID,
				Amount:              plan.Amount,
				Created:             plan.Created,
				Currency:            stripeToPbCurrency(plan.Currency),
				Interval:            stripeToPbInterval(plan.Interval),
				IntervalCount:       plan.IntervalCount,
				Livemode:            plan.Live,
				Metadata:            plan.Meta,
				Name:                plan.Name,
				StatementDescriptor: plan.Statement,
				TrialPeriodDays:     plan.TrialPeriod,
			},
		},
	}
}

// convert an error response from Stripe to a PlanResponse (error)
func respToPlanError(err *stripe.Error) *pb.PlanResponse {
	return &pb.PlanResponse{
		Responses: &pb.PlanResponse_Error{
			Error: respToError(err),
		},
	}
}

// convert a delete success response from Stripe to a DeletePlanResponse
func respToDeleteSuccess(plan *stripe.Plan) *pb.DeletePlanResponse {
	return &pb.DeletePlanResponse{
		Responses: &pb.DeletePlanResponse_Success{
			Success: &pb.DeletePlanSuccess{
				Id:      plan.ID,
				Deleted: plan.Deleted,
			},
		},
	}
}

// convert a delete error response from Stripe to a DeletePlanResponse
func respToDeleteError(err *stripe.Error) *pb.DeletePlanResponse {
	return &pb.DeletePlanResponse{
		Responses: &pb.DeletePlanResponse_Error{
			Error: respToError(err),
		},
	}
}

// constant conversions from stripe to protobuf - currency
func stripeToPbCurrency(c stripe.Currency) pb.Currency {
	return pb.Currency(pb.Currency_value[strings.ToUpper(fmt.Sprintf("%s", c))])
}

// constant conversions from stripe to protobuf - plan interval
func stripeToPbInterval(i stripe.PlanInterval) pb.Interval {
	return pb.Interval(pb.Interval_value[strings.Title(fmt.Sprintf("%s", i))])
}

// constant conversions from protobuf to stripe - currency
func pbToStripeCurrency(e pb.Currency) stripe.Currency {
	return stripe.Currency(strings.ToLower(pb.Currency_name[int32(e)]))
}

// constant conversions from protobuf to stripe - interval
func pbToStripeInterval(p pb.Interval) stripe.PlanInterval {
	return stripe.PlanInterval(strings.ToLower(pb.Interval_name[int32(p)]))
}
