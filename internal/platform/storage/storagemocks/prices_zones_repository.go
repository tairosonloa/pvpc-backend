// Code generated by mockery v2.32.4. DO NOT EDIT.

package storagemocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	pvpc "go-pvpc/internal"
)

// PricesZonesRepository is an autogenerated mock type for the PricesZonesRepository type
type PricesZonesRepository struct {
	mock.Mock
}

type PricesZonesRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *PricesZonesRepository) EXPECT() *PricesZonesRepository_Expecter {
	return &PricesZonesRepository_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: ctx
func (_m *PricesZonesRepository) GetAll(ctx context.Context) ([]pvpc.PricesZone, error) {
	ret := _m.Called(ctx)

	var r0 []pvpc.PricesZone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]pvpc.PricesZone, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []pvpc.PricesZone); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pvpc.PricesZone)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PricesZonesRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type PricesZonesRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PricesZonesRepository_Expecter) GetAll(ctx interface{}) *PricesZonesRepository_GetAll_Call {
	return &PricesZonesRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *PricesZonesRepository_GetAll_Call) Run(run func(ctx context.Context)) *PricesZonesRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *PricesZonesRepository_GetAll_Call) Return(_a0 []pvpc.PricesZone, _a1 error) *PricesZonesRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PricesZonesRepository_GetAll_Call) RunAndReturn(run func(context.Context) ([]pvpc.PricesZone, error)) *PricesZonesRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByExternalID provides a mock function with given fields: ctx, externalID
func (_m *PricesZonesRepository) GetByExternalID(ctx context.Context, externalID string) (pvpc.PricesZone, error) {
	ret := _m.Called(ctx, externalID)

	var r0 pvpc.PricesZone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (pvpc.PricesZone, error)); ok {
		return rf(ctx, externalID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) pvpc.PricesZone); ok {
		r0 = rf(ctx, externalID)
	} else {
		r0 = ret.Get(0).(pvpc.PricesZone)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, externalID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PricesZonesRepository_GetByExternalID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByExternalID'
type PricesZonesRepository_GetByExternalID_Call struct {
	*mock.Call
}

// GetByExternalID is a helper method to define mock.On call
//   - ctx context.Context
//   - externalID string
func (_e *PricesZonesRepository_Expecter) GetByExternalID(ctx interface{}, externalID interface{}) *PricesZonesRepository_GetByExternalID_Call {
	return &PricesZonesRepository_GetByExternalID_Call{Call: _e.mock.On("GetByExternalID", ctx, externalID)}
}

func (_c *PricesZonesRepository_GetByExternalID_Call) Run(run func(ctx context.Context, externalID string)) *PricesZonesRepository_GetByExternalID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *PricesZonesRepository_GetByExternalID_Call) Return(_a0 pvpc.PricesZone, _a1 error) *PricesZonesRepository_GetByExternalID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PricesZonesRepository_GetByExternalID_Call) RunAndReturn(run func(context.Context, string) (pvpc.PricesZone, error)) *PricesZonesRepository_GetByExternalID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *PricesZonesRepository) GetByID(ctx context.Context, id pvpc.PricesZoneID) (pvpc.PricesZone, error) {
	ret := _m.Called(ctx, id)

	var r0 pvpc.PricesZone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pvpc.PricesZoneID) (pvpc.PricesZone, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pvpc.PricesZoneID) pvpc.PricesZone); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(pvpc.PricesZone)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pvpc.PricesZoneID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PricesZonesRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type PricesZonesRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id pvpc.PricesZoneID
func (_e *PricesZonesRepository_Expecter) GetByID(ctx interface{}, id interface{}) *PricesZonesRepository_GetByID_Call {
	return &PricesZonesRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *PricesZonesRepository_GetByID_Call) Run(run func(ctx context.Context, id pvpc.PricesZoneID)) *PricesZonesRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pvpc.PricesZoneID))
	})
	return _c
}

func (_c *PricesZonesRepository_GetByID_Call) Return(_a0 pvpc.PricesZone, _a1 error) *PricesZonesRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PricesZonesRepository_GetByID_Call) RunAndReturn(run func(context.Context, pvpc.PricesZoneID) (pvpc.PricesZone, error)) *PricesZonesRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewPricesZonesRepository creates a new instance of PricesZonesRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPricesZonesRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PricesZonesRepository {
	mock := &PricesZonesRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}