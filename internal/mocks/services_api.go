// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mgdb "github.com/acrosdale/gomongo/internal/db/mgdb"
	mock "github.com/stretchr/testify/mock"
)

// ApiServicesMock is an autogenerated mock type for the ApiServiceInterface type
type ApiServicesMock struct {
	mock.Mock
}

// DeleteProduct provides a mock function with given fields: ctx, filters
func (_m *ApiServicesMock) DeleteProduct(ctx context.Context, filters map[string]interface{}) (int64, error) {
	ret := _m.Called(ctx, filters)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) int64); ok {
		r0 = rf(ctx, filters)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProduct provides a mock function with given fields: ctx, filters
func (_m *ApiServicesMock) GetProduct(ctx context.Context, filters map[string]interface{}) (mgdb.Product, error) {
	ret := _m.Called(ctx, filters)

	var r0 mgdb.Product
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}) mgdb.Product); ok {
		r0 = rf(ctx, filters)
	} else {
		r0 = ret.Get(0).(mgdb.Product)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}) error); ok {
		r1 = rf(ctx, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertProduct provides a mock function with given fields: ctx, document
func (_m *ApiServicesMock) InsertProduct(ctx context.Context, document mgdb.Product) (string, error) {
	ret := _m.Called(ctx, document)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, mgdb.Product) string); ok {
		r0 = rf(ctx, document)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, mgdb.Product) error); ok {
		r1 = rf(ctx, document)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProduct provides a mock function with given fields: ctx, filters, document
func (_m *ApiServicesMock) UpdateProduct(ctx context.Context, filters map[string]interface{}, document mgdb.Product) (int64, error) {
	ret := _m.Called(ctx, filters, document)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, mgdb.Product) int64); ok {
		r0 = rf(ctx, filters, document)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, mgdb.Product) error); ok {
		r1 = rf(ctx, filters, document)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}