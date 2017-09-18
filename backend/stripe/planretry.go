package stripe

import (
	"github.com/BTBurke/recur/pb"
	"github.com/cenkalti/backoff"

	"github.com/stripe/stripe-go"
)

type planAction int

const (
	planCreate planAction = iota
	planUpdate
	planDelete
	planGet
)

func retryablePlan(params *stripe.PlanParams, api planClient, p *pb.PlanResponse, action planAction) backoff.Operation {
	return func() error {
		var plan = new(stripe.Plan)
		var err error
		switch action {
		case planCreate:
			plan, err = api.New(params)
		case planUpdate:
			plan, err = api.Update(params.ID, params)
		case planGet:
			plan, err = api.Get(params.ID, params)
		default:
		}
		if err != nil {
			switch err.(type) {
			case *stripe.Error:
				*p = *respToPlanError(err.(*stripe.Error))
				return nil
			default:
				return err
			}
		}
		*p = *respToPlanSuccess(plan)
		return nil
	}
}

func retryablePlanDelete(params *stripe.PlanParams, api planClient, p *pb.DeletePlanResponse, action planAction) backoff.Operation {
	return func() error {
		var plan = new(stripe.Plan)
		var err error
		switch action {
		case planDelete:
			plan, err = api.Del(params.ID, params)
		default:
		}
		if err != nil {
			switch err.(type) {
			case *stripe.Error:
				*p = *respToDeleteError(err.(*stripe.Error))
				return nil
			default:
				return err
			}
		}
		*p = *respToDeleteSuccess(plan)
		return nil
	}
}

func retryablePlanList(params *stripe.PlanListParams, api planClient, p *planStreamer) backoff.Operation {
	return func() error {
		p.iter = api.List(params)
		return nil
	}
}
