package stripe

import (
	"os"
	"strings"
	"testing"

	"github.com/BTBurke/recur/pb"
	"github.com/stretchr/testify/assert"
	context "golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

func getAPIKeyOrSkip(t *testing.T) string {
	key := os.Getenv("STRIPE_KEY")
	if len(key) == 0 || !strings.Contains(key, "sk_test") {
		t.Skipf("No valid stripe testing key found, skipping Stripe integration tests. Set env STRIPE_KEY to run tests.")
		return ""
	}

	return key
}

func TestPlanIntegration(t *testing.T) {
	key := getAPIKeyOrSkip(t)
	client := NewPlanClient(key, log.New())

	if err := deleteAllExistingPlans(client); err != nil {
		t.Fatalf("Failed to delete existing plans before starting integration test: %s", err)
	}

	r1 := &pb.CreatePlanRequest{
		Id:                  "test-plan-5",
		Amount:              1000,
		Currency:            pb.Currency_USD,
		Name:                "test",
		Interval:            pb.Interval_Month,
		IntervalCount:       1,
		Metadata:            map[string]string{"test": "test1"},
		StatementDescriptor: "nothing",
		TrialPeriodDays:     1,
	}
	r2 := &pb.CreatePlanRequest{
		Id:                  "test-plan-6",
		Amount:              999,
		Currency:            pb.Currency_AED,
		Name:                "test2",
		Interval:            pb.Interval_Month,
		IntervalCount:       2,
		Metadata:            map[string]string{"test": "test2"},
		StatementDescriptor: "nothing2",
		TrialPeriodDays:     2,
	}
	r3 := &pb.CreatePlanRequest{
		Id: "test-plan-x",
	}
	tt := []struct {
		Name      string
		Reqs      []*pb.CreatePlanRequest
		ShouldErr bool
	}{
		{Name: "add multiple plans", Reqs: []*pb.CreatePlanRequest{r1, r2}, ShouldErr: false},
		{Name: "fail validation", Reqs: []*pb.CreatePlanRequest{r3}, ShouldErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {

			// create plans
			var expPlans []*pb.Plan
			for _, req := range tc.Reqs {
				plan, err := client.Create(context.Background(), req)
				switch tc.ShouldErr {
				case true:
					assert.Error(t, err)
					return
				default:
					assert.NoError(t, err)
					got := plan.GetSuccess()
					assert.Equal(t, req.Id, got.Id, "id")
					assert.Equal(t, req.Amount, got.Amount, "amount")
					assert.Equal(t, req.Currency, got.Currency, "currency")
					assert.Equal(t, req.Name, got.Name, "name")
					assert.Equal(t, req.Interval, got.Interval, "interval")
					assert.Equal(t, req.IntervalCount, got.IntervalCount, "interval count")
					assert.Equal(t, req.Metadata, got.Metadata, "metadata")
					assert.Equal(t, req.StatementDescriptor, got.StatementDescriptor, "statement descriptor")
					assert.Equal(t, req.TrialPeriodDays, got.TrialPeriodDays, "trial period")
					expPlans = append([]*pb.Plan{got}, expPlans...)
				}
			}

			// Check all plans will list
			plans, err := client.List(context.Background(), nil)
			if err != nil {
				t.Fatalf("Unexpected error listing plans: %s", err)
			}
			var gotPlans []*pb.Plan
			for plans.Next() {
				planResp := plans.Current()
				gotPlans = append(gotPlans, planResp.GetSuccess())
			}
			for i, expPlan := range expPlans {
				assert.Equal(t, expPlan, gotPlans[i])
			}

			// update single plan
			resp, err := client.Update(context.Background(), &pb.UpdatePlanRequest{
				Id:   r1.Id,
				Name: "updated name",
			})
			assert.NoError(t, err)
			assert.Equal(t, "updated name", resp.GetSuccess().Name)

			// get single plan
			resp2, err := client.Get(context.Background(), &pb.GetPlanRequest{
				Id: r2.Id,
			})
			assert.NoError(t, err)
			assert.Equal(t, expPlans[0], resp2.GetSuccess())
		})
	}
}

func deleteAllExistingPlans(client *StripePlanClient) error {
	plans, err := client.List(context.Background(), nil)
	if err != nil {
		return err
	}
	for plans.Next() {
		planResp := plans.Current()
		plan := planResp.GetSuccess()
		delResp, err := client.Delete(context.Background(), &pb.DeletePlanRequest{
			Id: plan.Id,
		})
		if err != nil {
			return err
		}
		log.Printf("Deleted: %#v\n", delResp.GetSuccess().Id)
	}
	return nil
}
