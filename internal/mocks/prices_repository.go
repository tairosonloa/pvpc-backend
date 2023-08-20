// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	"pvpc-backend/internal/domain"
)

// PricesRepository is an autogenerated mock type for the PricesRepository type
type PricesRepository struct {
	mock.Mock
}

type PricesRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *PricesRepository) EXPECT() *PricesRepository_Expecter {
	return &PricesRepository_Expecter{mock: &_m.Mock}
}

// Save provides a mock function with given fields: ctx, prices
func (_m *PricesRepository) Save(ctx context.Context, prices []domain.Prices) error {
	ret := _m.Called(ctx, prices)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []domain.Prices) error); ok {
		r0 = rf(ctx, prices)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PricesRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type PricesRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - ctx context.Context
//   - prices []domain.Prices
func (_e *PricesRepository_Expecter) Save(ctx interface{}, prices interface{}) *PricesRepository_Save_Call {
	return &PricesRepository_Save_Call{Call: _e.mock.On("Save", ctx, prices)}
}

func (_c *PricesRepository_Save_Call) Run(run func(ctx context.Context, prices []domain.Prices)) *PricesRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]domain.Prices))
	})
	return _c
}

func (_c *PricesRepository_Save_Call) Return(_a0 error) *PricesRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PricesRepository_Save_Call) RunAndReturn(run func(context.Context, []domain.Prices) error) *PricesRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewPricesRepository creates a new instance of PricesRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPricesRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PricesRepository {
	mock := &PricesRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}