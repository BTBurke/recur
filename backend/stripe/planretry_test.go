package stripe

import (
	"fmt"
	"testing"
	"time"

	context "golang.org/x/net/context"

	"github.com/BTBurke/recur/pb"
	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/plan"
)

type mockPlan struct {
	delay  time.Duration
	called int
	mock.Mock
}

func (m *mockPlan) New(params *stripe.PlanParams) (*stripe.Plan, error) {
	args := m.Called()
	time.Sleep(m.delay)
	if m.called > 0 && m.delay == 0 {
		return args.Get(0).(*stripe.Plan), args.Error(1)
	}
	m.called++
	return nil, fmt.Errorf("test retry")
}
func (m *mockPlan) Get(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	args := m.Called()
	time.Sleep(m.delay)
	if m.called > 0 && m.delay == 0 {
		return args.Get(0).(*stripe.Plan), args.Error(1)
	}
	m.called++
	return nil, fmt.Errorf("test retry")
}
func (m *mockPlan) Update(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	args := m.Called()
	time.Sleep(m.delay)
	if m.called > 0 && m.delay == 0 {
		return args.Get(0).(*stripe.Plan), args.Error(1)
	}
	m.called++
	return nil, fmt.Errorf("test retry")
}
func (m *mockPlan) Del(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	args := m.Called()
	time.Sleep(m.delay)
	if m.called > 0 {
		return args.Get(0).(*stripe.Plan), args.Error(1)
	}
	m.called++
	return nil, fmt.Errorf("test retry")
}
func (m *mockPlan) List(params *stripe.PlanListParams) *plan.Iter {
	m.called++
	args := m.Called(params)
	time.Sleep(m.delay)
	return args.Get(0).(*plan.Iter)
}

func TestRetryablePlan(t *testing.T) {
	pln := &stripe.Plan{
		ID:       "test",
		Live:     true,
		Amount:   1000,
		Currency: stripe.Currency("usd"),
		Interval: stripe.PlanInterval("month"),
	}
	tt := []struct {
		Name      string
		Method    string
		Action    planAction
		Timeout   time.Duration
		Delay     time.Duration
		ShouldErr bool
	}{
		{Name: "create no timeout", Method: "New", Action: planCreate, Timeout: 1 * time.Second, ShouldErr: false},
		{Name: "get no timeout", Method: "Get", Action: planGet, Timeout: 1 * time.Second, ShouldErr: false},
		{Name: "update no timeout", Method: "Update", Action: planUpdate, Timeout: 1 * time.Second, ShouldErr: false},
		{Name: "update times out", Method: "Update", Action: planUpdate, Timeout: 100 * time.Millisecond, Delay: 120 * time.Millisecond, ShouldErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			mck := &mockPlan{delay: tc.Delay}
			resp := new(pb.PlanResponse)
			ctx, cancelFunc := context.WithTimeout(context.Background(), tc.Timeout)
			defer cancelFunc()

			mck.On(tc.Method).Return(pln, nil)
			err := backoff.Retry(
				retryablePlan(&stripe.PlanParams{}, mck, resp, tc.Action),
				backoff.WithContext(backoff.NewExponentialBackOff(), ctx),
			)
			switch tc.ShouldErr {
			case true:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
				mck.AssertNumberOfCalls(t, tc.Method, 2)
			}

		})
	}
}
