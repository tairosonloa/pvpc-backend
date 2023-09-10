// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	domain "pvpc-backend/internal/domain"
)

// ZonesRepository is an autogenerated mock type for the ZonesRepository type
type ZonesRepository struct {
	mock.Mock
}

type ZonesRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *ZonesRepository) EXPECT() *ZonesRepository_Expecter {
	return &ZonesRepository_Expecter{mock: &_m.Mock}
}

// GetAll provides a mock function with given fields: ctx
func (_m *ZonesRepository) GetAll(ctx context.Context) ([]domain.Zone, error) {
	ret := _m.Called(ctx)

	var r0 []domain.Zone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]domain.Zone, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Zone); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Zone)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ZonesRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type ZonesRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ZonesRepository_Expecter) GetAll(ctx interface{}) *ZonesRepository_GetAll_Call {
	return &ZonesRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *ZonesRepository_GetAll_Call) Run(run func(ctx context.Context)) *ZonesRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ZonesRepository_GetAll_Call) Return(_a0 []domain.Zone, _a1 error) *ZonesRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ZonesRepository_GetAll_Call) RunAndReturn(run func(context.Context) ([]domain.Zone, error)) *ZonesRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByExternalID provides a mock function with given fields: ctx, externalID
func (_m *ZonesRepository) GetByExternalID(ctx context.Context, externalID string) (domain.Zone, error) {
	ret := _m.Called(ctx, externalID)

	var r0 domain.Zone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.Zone, error)); ok {
		return rf(ctx, externalID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Zone); ok {
		r0 = rf(ctx, externalID)
	} else {
		r0 = ret.Get(0).(domain.Zone)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, externalID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ZonesRepository_GetByExternalID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByExternalID'
type ZonesRepository_GetByExternalID_Call struct {
	*mock.Call
}

// GetByExternalID is a helper method to define mock.On call
//   - ctx context.Context
//   - externalID string
func (_e *ZonesRepository_Expecter) GetByExternalID(ctx interface{}, externalID interface{}) *ZonesRepository_GetByExternalID_Call {
	return &ZonesRepository_GetByExternalID_Call{Call: _e.mock.On("GetByExternalID", ctx, externalID)}
}

func (_c *ZonesRepository_GetByExternalID_Call) Run(run func(ctx context.Context, externalID string)) *ZonesRepository_GetByExternalID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ZonesRepository_GetByExternalID_Call) Return(_a0 domain.Zone, _a1 error) *ZonesRepository_GetByExternalID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ZonesRepository_GetByExternalID_Call) RunAndReturn(run func(context.Context, string) (domain.Zone, error)) *ZonesRepository_GetByExternalID_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *ZonesRepository) GetByID(ctx context.Context, id domain.ZoneID) (domain.Zone, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Zone
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ZoneID) (domain.Zone, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.ZoneID) domain.Zone); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Zone)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.ZoneID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ZonesRepository_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type ZonesRepository_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id domain.ZoneID
func (_e *ZonesRepository_Expecter) GetByID(ctx interface{}, id interface{}) *ZonesRepository_GetByID_Call {
	return &ZonesRepository_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *ZonesRepository_GetByID_Call) Run(run func(ctx context.Context, id domain.ZoneID)) *ZonesRepository_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ZoneID))
	})
	return _c
}

func (_c *ZonesRepository_GetByID_Call) Return(_a0 domain.Zone, _a1 error) *ZonesRepository_GetByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ZonesRepository_GetByID_Call) RunAndReturn(run func(context.Context, domain.ZoneID) (domain.Zone, error)) *ZonesRepository_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// NewZonesRepository creates a new instance of ZonesRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewZonesRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ZonesRepository {
	mock := &ZonesRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
