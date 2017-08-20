package backend

import (
	"github.com/BTBurke/recur/pb"
	context "golang.org/x/net/context"
)

// GenericStreamer is a subset of methods to allow streaming responses from the backend
type GenericStreamer interface {
	Next() bool
	Current() interface{}
}

// PlanClient is an interface for actions related to CRUD operations on a backend (e.g. Stripe)
type PlanClient interface {
	Create(ctx context.Context, req *pb.CreatePlanRequest) (*pb.PlanResponse, error)
	// Update(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.PlanResponse, error)
	// Delete(ctx context.Context, req *pb.DeletePlanRequest) (*pb.DeletePlanResponse, error)
	// Get(ctx context.Context, req *pb.GetPlanRequest) (*pb.PlanResponse, error)
	// List(ctx context.Context, req *pb.ListRequest) (*GenericStreamer, error)
}
