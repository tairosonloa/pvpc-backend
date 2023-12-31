// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"

	domain "pvpc-backend/internal/domain"
)

// PricesProvider is an autogenerated mock type for the PricesProvider type
type PricesProvider struct {
	mock.Mock
}

type PricesProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *PricesProvider) EXPECT() *PricesProvider_Expecter {
	return &PricesProvider_Expecter{mock: &_m.Mock}
}

// FetchPVPCPrices provides a mock function with given fields: ctx, zones, date
func (_m *PricesProvider) FetchPVPCPrices(ctx context.Context, zones []domain.Zone, date time.Time) ([]domain.Prices, error) {
	ret := _m.Called(ctx, zones, date)

	var r0 []domain.Prices
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.Zone, time.Time) ([]domain.Prices, error)); ok {
		return rf(ctx, zones, date)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []domain.Zone, time.Time) []domain.Prices); ok {
		r0 = rf(ctx, zones, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Prices)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []domain.Zone, time.Time) error); ok {
		r1 = rf(ctx, zones, date)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PricesProvider_FetchPVPCPrices_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchPVPCPrices'
type PricesProvider_FetchPVPCPrices_Call struct {
	*mock.Call
}

// FetchPVPCPrices is a helper method to define mock.On call
//   - ctx context.Context
//   - zones []domain.Zone
//   - date time.Time
func (_e *PricesProvider_Expecter) FetchPVPCPrices(ctx interface{}, zones interface{}, date interface{}) *PricesProvider_FetchPVPCPrices_Call {
	return &PricesProvider_FetchPVPCPrices_Call{Call: _e.mock.On("FetchPVPCPrices", ctx, zones, date)}
}

func (_c *PricesProvider_FetchPVPCPrices_Call) Run(run func(ctx context.Context, zones []domain.Zone, date time.Time)) *PricesProvider_FetchPVPCPrices_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]domain.Zone), args[2].(time.Time))
	})
	return _c
}

func (_c *PricesProvider_FetchPVPCPrices_Call) Return(_a0 []domain.Prices, _a1 error) *PricesProvider_FetchPVPCPrices_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PricesProvider_FetchPVPCPrices_Call) RunAndReturn(run func(context.Context, []domain.Zone, time.Time) ([]domain.Prices, error)) *PricesProvider_FetchPVPCPrices_Call {
	_c.Call.Return(run)
	return _c
}

// NewPricesProvider creates a new instance of PricesProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPricesProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *PricesProvider {
	mock := &PricesProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
