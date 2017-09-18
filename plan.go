package recur

import (
	"google.golang.org/grpc"

	"github.com/BTBurke/recur/backend"
	"github.com/BTBurke/recur/pb"
	context "golang.org/x/net/context"
)

type PlanClient struct {
	backend backend.PlanClient
}

// CreatePlan is the GRPC endpoint to create a plan.
func (c *PlanClient) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest, opts ...grpc.CallOption) (*pb.PlanResponse, error) {
	return nil, nil
}

// Create creates a plan with a default context
func (c *PlanClient) Create(req *pb.CreatePlanRequest) (*pb.PlanResponse, error) {
	return nil, nil
}

// CreateWithCtx creates a plan with a custom context
func (c *PlanClient) CreateWithCtx(ctx context.Context, req *pb.CreatePlanRequest) (*pb.PlanResponse, error) {
	return nil, nil
}

func (c *PlanClient) create(ctx context.Context, req *pb.CreatePlanRequest) (*pb.PlanResponse, error) {
	return nil, nil
}
