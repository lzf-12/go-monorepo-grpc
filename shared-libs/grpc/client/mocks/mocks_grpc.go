// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package grpc

import (
	"context"
	"pb_schemas/inventory/v1"
	"pb_schemas/user/v1"

	mock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// NewMockInvClient creates a new instance of MockInvClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInvClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInvClient {
	mock := &MockInvClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockInvClient is an autogenerated mock type for the InvClient type
type MockInvClient struct {
	mock.Mock
}

type MockInvClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInvClient) EXPECT() *MockInvClient_Expecter {
	return &MockInvClient_Expecter{mock: &_m.Mock}
}

// CheckStock provides a mock function for the type MockInvClient
func (_mock *MockInvClient) CheckStock(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption) (*inventoryv1.InventoryStatusResponse, error) {
	var tmpRet mock.Arguments
	if len(opts) > 0 {
		tmpRet = _mock.Called(ctx, in, opts)
	} else {
		tmpRet = _mock.Called(ctx, in)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for CheckStock")
	}

	var r0 *inventoryv1.InventoryStatusResponse
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) (*inventoryv1.InventoryStatusResponse, error)); ok {
		return returnFunc(ctx, in, opts...)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) *inventoryv1.InventoryStatusResponse); ok {
		r0 = returnFunc(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*inventoryv1.InventoryStatusResponse)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) error); ok {
		r1 = returnFunc(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockInvClient_CheckStock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckStock'
type MockInvClient_CheckStock_Call struct {
	*mock.Call
}

// CheckStock is a helper method to define mock.On call
//   - ctx context.Context
//   - in *inventoryv1.StandardInventoryRequest
//   - opts ...grpc.CallOption
func (_e *MockInvClient_Expecter) CheckStock(ctx interface{}, in interface{}, opts ...interface{}) *MockInvClient_CheckStock_Call {
	return &MockInvClient_CheckStock_Call{Call: _e.mock.On("CheckStock",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockInvClient_CheckStock_Call) Run(run func(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption)) *MockInvClient_CheckStock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *inventoryv1.StandardInventoryRequest
		if args[1] != nil {
			arg1 = args[1].(*inventoryv1.StandardInventoryRequest)
		}
		var arg2 []grpc.CallOption
		var variadicArgs []grpc.CallOption
		if len(args) > 2 {
			variadicArgs = args[2].([]grpc.CallOption)
		}
		arg2 = variadicArgs
		run(
			arg0,
			arg1,
			arg2...,
		)
	})
	return _c
}

func (_c *MockInvClient_CheckStock_Call) Return(inventoryStatusResponse *inventoryv1.InventoryStatusResponse, err error) *MockInvClient_CheckStock_Call {
	_c.Call.Return(inventoryStatusResponse, err)
	return _c
}

func (_c *MockInvClient_CheckStock_Call) RunAndReturn(run func(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption) (*inventoryv1.InventoryStatusResponse, error)) *MockInvClient_CheckStock_Call {
	_c.Call.Return(run)
	return _c
}

// ReleaseStock provides a mock function for the type MockInvClient
func (_mock *MockInvClient) ReleaseStock(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption) (*inventoryv1.InventoryReservationResponse, error) {
	var tmpRet mock.Arguments
	if len(opts) > 0 {
		tmpRet = _mock.Called(ctx, in, opts)
	} else {
		tmpRet = _mock.Called(ctx, in)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for ReleaseStock")
	}

	var r0 *inventoryv1.InventoryReservationResponse
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) (*inventoryv1.InventoryReservationResponse, error)); ok {
		return returnFunc(ctx, in, opts...)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) *inventoryv1.InventoryReservationResponse); ok {
		r0 = returnFunc(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*inventoryv1.InventoryReservationResponse)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) error); ok {
		r1 = returnFunc(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockInvClient_ReleaseStock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReleaseStock'
type MockInvClient_ReleaseStock_Call struct {
	*mock.Call
}

// ReleaseStock is a helper method to define mock.On call
//   - ctx context.Context
//   - in *inventoryv1.StandardInventoryRequest
//   - opts ...grpc.CallOption
func (_e *MockInvClient_Expecter) ReleaseStock(ctx interface{}, in interface{}, opts ...interface{}) *MockInvClient_ReleaseStock_Call {
	return &MockInvClient_ReleaseStock_Call{Call: _e.mock.On("ReleaseStock",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockInvClient_ReleaseStock_Call) Run(run func(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption)) *MockInvClient_ReleaseStock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *inventoryv1.StandardInventoryRequest
		if args[1] != nil {
			arg1 = args[1].(*inventoryv1.StandardInventoryRequest)
		}
		var arg2 []grpc.CallOption
		var variadicArgs []grpc.CallOption
		if len(args) > 2 {
			variadicArgs = args[2].([]grpc.CallOption)
		}
		arg2 = variadicArgs
		run(
			arg0,
			arg1,
			arg2...,
		)
	})
	return _c
}

func (_c *MockInvClient_ReleaseStock_Call) Return(inventoryReservationResponse *inventoryv1.InventoryReservationResponse, err error) *MockInvClient_ReleaseStock_Call {
	_c.Call.Return(inventoryReservationResponse, err)
	return _c
}

func (_c *MockInvClient_ReleaseStock_Call) RunAndReturn(run func(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption) (*inventoryv1.InventoryReservationResponse, error)) *MockInvClient_ReleaseStock_Call {
	_c.Call.Return(run)
	return _c
}

// ReserveStock provides a mock function for the type MockInvClient
func (_mock *MockInvClient) ReserveStock(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption) (*inventoryv1.InventoryReservationResponse, error) {
	var tmpRet mock.Arguments
	if len(opts) > 0 {
		tmpRet = _mock.Called(ctx, in, opts)
	} else {
		tmpRet = _mock.Called(ctx, in)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for ReserveStock")
	}

	var r0 *inventoryv1.InventoryReservationResponse
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) (*inventoryv1.InventoryReservationResponse, error)); ok {
		return returnFunc(ctx, in, opts...)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) *inventoryv1.InventoryReservationResponse); ok {
		r0 = returnFunc(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*inventoryv1.InventoryReservationResponse)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *inventoryv1.StandardInventoryRequest, ...grpc.CallOption) error); ok {
		r1 = returnFunc(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockInvClient_ReserveStock_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReserveStock'
type MockInvClient_ReserveStock_Call struct {
	*mock.Call
}

// ReserveStock is a helper method to define mock.On call
//   - ctx context.Context
//   - in *inventoryv1.StandardInventoryRequest
//   - opts ...grpc.CallOption
func (_e *MockInvClient_Expecter) ReserveStock(ctx interface{}, in interface{}, opts ...interface{}) *MockInvClient_ReserveStock_Call {
	return &MockInvClient_ReserveStock_Call{Call: _e.mock.On("ReserveStock",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockInvClient_ReserveStock_Call) Run(run func(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption)) *MockInvClient_ReserveStock_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *inventoryv1.StandardInventoryRequest
		if args[1] != nil {
			arg1 = args[1].(*inventoryv1.StandardInventoryRequest)
		}
		var arg2 []grpc.CallOption
		var variadicArgs []grpc.CallOption
		if len(args) > 2 {
			variadicArgs = args[2].([]grpc.CallOption)
		}
		arg2 = variadicArgs
		run(
			arg0,
			arg1,
			arg2...,
		)
	})
	return _c
}

func (_c *MockInvClient_ReserveStock_Call) Return(inventoryReservationResponse *inventoryv1.InventoryReservationResponse, err error) *MockInvClient_ReserveStock_Call {
	_c.Call.Return(inventoryReservationResponse, err)
	return _c
}

func (_c *MockInvClient_ReserveStock_Call) RunAndReturn(run func(ctx context.Context, in *inventoryv1.StandardInventoryRequest, opts ...grpc.CallOption) (*inventoryv1.InventoryReservationResponse, error)) *MockInvClient_ReserveStock_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserClient creates a new instance of MockUserClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserClient {
	mock := &MockUserClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockUserClient is an autogenerated mock type for the UserClient type
type MockUserClient struct {
	mock.Mock
}

type MockUserClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserClient) EXPECT() *MockUserClient_Expecter {
	return &MockUserClient_Expecter{mock: &_m.Mock}
}

// ValidateToken provides a mock function for the type MockUserClient
func (_mock *MockUserClient) ValidateToken(ctx context.Context, in *userv1.ValidateTokenRequest, opts ...grpc.CallOption) (*userv1.ValidateTokenResponse, error) {
	var tmpRet mock.Arguments
	if len(opts) > 0 {
		tmpRet = _mock.Called(ctx, in, opts)
	} else {
		tmpRet = _mock.Called(ctx, in)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *userv1.ValidateTokenResponse
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *userv1.ValidateTokenRequest, ...grpc.CallOption) (*userv1.ValidateTokenResponse, error)); ok {
		return returnFunc(ctx, in, opts...)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *userv1.ValidateTokenRequest, ...grpc.CallOption) *userv1.ValidateTokenResponse); ok {
		r0 = returnFunc(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userv1.ValidateTokenResponse)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *userv1.ValidateTokenRequest, ...grpc.CallOption) error); ok {
		r1 = returnFunc(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockUserClient_ValidateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateToken'
type MockUserClient_ValidateToken_Call struct {
	*mock.Call
}

// ValidateToken is a helper method to define mock.On call
//   - ctx context.Context
//   - in *userv1.ValidateTokenRequest
//   - opts ...grpc.CallOption
func (_e *MockUserClient_Expecter) ValidateToken(ctx interface{}, in interface{}, opts ...interface{}) *MockUserClient_ValidateToken_Call {
	return &MockUserClient_ValidateToken_Call{Call: _e.mock.On("ValidateToken",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockUserClient_ValidateToken_Call) Run(run func(ctx context.Context, in *userv1.ValidateTokenRequest, opts ...grpc.CallOption)) *MockUserClient_ValidateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *userv1.ValidateTokenRequest
		if args[1] != nil {
			arg1 = args[1].(*userv1.ValidateTokenRequest)
		}
		var arg2 []grpc.CallOption
		var variadicArgs []grpc.CallOption
		if len(args) > 2 {
			variadicArgs = args[2].([]grpc.CallOption)
		}
		arg2 = variadicArgs
		run(
			arg0,
			arg1,
			arg2...,
		)
	})
	return _c
}

func (_c *MockUserClient_ValidateToken_Call) Return(validateTokenResponse *userv1.ValidateTokenResponse, err error) *MockUserClient_ValidateToken_Call {
	_c.Call.Return(validateTokenResponse, err)
	return _c
}

func (_c *MockUserClient_ValidateToken_Call) RunAndReturn(run func(ctx context.Context, in *userv1.ValidateTokenRequest, opts ...grpc.CallOption) (*userv1.ValidateTokenResponse, error)) *MockUserClient_ValidateToken_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockUserGrpcClientInterface creates a new instance of MockUserGrpcClientInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUserGrpcClientInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUserGrpcClientInterface {
	mock := &MockUserGrpcClientInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockUserGrpcClientInterface is an autogenerated mock type for the UserGrpcClientInterface type
type MockUserGrpcClientInterface struct {
	mock.Mock
}

type MockUserGrpcClientInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUserGrpcClientInterface) EXPECT() *MockUserGrpcClientInterface_Expecter {
	return &MockUserGrpcClientInterface_Expecter{mock: &_m.Mock}
}

// ValidateToken provides a mock function for the type MockUserGrpcClientInterface
func (_mock *MockUserGrpcClientInterface) ValidateToken(ctx context.Context, in *userv1.ValidateTokenRequest, opts ...grpc.CallOption) (*userv1.ValidateTokenResponse, error) {
	var tmpRet mock.Arguments
	if len(opts) > 0 {
		tmpRet = _mock.Called(ctx, in, opts)
	} else {
		tmpRet = _mock.Called(ctx, in)
	}
	ret := tmpRet

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *userv1.ValidateTokenResponse
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, *userv1.ValidateTokenRequest, ...grpc.CallOption) (*userv1.ValidateTokenResponse, error)); ok {
		return returnFunc(ctx, in, opts...)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, *userv1.ValidateTokenRequest, ...grpc.CallOption) *userv1.ValidateTokenResponse); ok {
		r0 = returnFunc(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userv1.ValidateTokenResponse)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, *userv1.ValidateTokenRequest, ...grpc.CallOption) error); ok {
		r1 = returnFunc(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockUserGrpcClientInterface_ValidateToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateToken'
type MockUserGrpcClientInterface_ValidateToken_Call struct {
	*mock.Call
}

// ValidateToken is a helper method to define mock.On call
//   - ctx context.Context
//   - in *userv1.ValidateTokenRequest
//   - opts ...grpc.CallOption
func (_e *MockUserGrpcClientInterface_Expecter) ValidateToken(ctx interface{}, in interface{}, opts ...interface{}) *MockUserGrpcClientInterface_ValidateToken_Call {
	return &MockUserGrpcClientInterface_ValidateToken_Call{Call: _e.mock.On("ValidateToken",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockUserGrpcClientInterface_ValidateToken_Call) Run(run func(ctx context.Context, in *userv1.ValidateTokenRequest, opts ...grpc.CallOption)) *MockUserGrpcClientInterface_ValidateToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 *userv1.ValidateTokenRequest
		if args[1] != nil {
			arg1 = args[1].(*userv1.ValidateTokenRequest)
		}
		var arg2 []grpc.CallOption
		var variadicArgs []grpc.CallOption
		if len(args) > 2 {
			variadicArgs = args[2].([]grpc.CallOption)
		}
		arg2 = variadicArgs
		run(
			arg0,
			arg1,
			arg2...,
		)
	})
	return _c
}

func (_c *MockUserGrpcClientInterface_ValidateToken_Call) Return(validateTokenResponse *userv1.ValidateTokenResponse, err error) *MockUserGrpcClientInterface_ValidateToken_Call {
	_c.Call.Return(validateTokenResponse, err)
	return _c
}

func (_c *MockUserGrpcClientInterface_ValidateToken_Call) RunAndReturn(run func(ctx context.Context, in *userv1.ValidateTokenRequest, opts ...grpc.CallOption) (*userv1.ValidateTokenResponse, error)) *MockUserGrpcClientInterface_ValidateToken_Call {
	_c.Call.Return(run)
	return _c
}
