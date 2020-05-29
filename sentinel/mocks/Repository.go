// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	resource "github.com/bithippie/go-sentinel/sentinel/models/resource"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *Repository) Create(_a0 *resource.Input) (resource.Element, error) {
	ret := _m.Called(_a0)

	var r0 resource.Element
	if rf, ok := ret.Get(0).(func(*resource.Input) resource.Element); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(resource.Element)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*resource.Input) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: _a0
func (_m *Repository) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields:
func (_m *Repository) Get() (resource.Response, error) {
	ret := _m.Called()

	var r0 resource.Response
	if rf, ok := ret.Get(0).(func() resource.Response); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(resource.Response)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: _a0
func (_m *Repository) GetByID(_a0 string) (resource.Element, error) {
	ret := _m.Called(_a0)

	var r0 resource.Element
	if rf, ok := ret.Get(0).(func(string) resource.Element); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(resource.Element)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *Repository) Update(_a0 resource.Element, _a1 *resource.Input) (resource.Element, error) {
	ret := _m.Called(_a0, _a1)

	var r0 resource.Element
	if rf, ok := ret.Get(0).(func(resource.Element, *resource.Input) resource.Element); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(resource.Element)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(resource.Element, *resource.Input) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
