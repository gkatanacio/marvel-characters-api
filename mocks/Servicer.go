// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	marvel "github.com/gkatanacio/marvel-characters-api/internal/marvel"
	mock "github.com/stretchr/testify/mock"
)

// Servicer is an autogenerated mock type for the Servicer type
type Servicer struct {
	mock.Mock
}

// GetAllCharacterIds provides a mock function with given fields:
func (_m *Servicer) GetAllCharacterIds() ([]int, error) {
	ret := _m.Called()

	var r0 []int
	if rf, ok := ret.Get(0).(func() []int); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCharacter provides a mock function with given fields: id
func (_m *Servicer) GetCharacter(id int) (*marvel.Character, error) {
	ret := _m.Called(id)

	var r0 *marvel.Character
	if rf, ok := ret.Get(0).(func(int) *marvel.Character); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*marvel.Character)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReloadCache provides a mock function with given fields:
func (_m *Servicer) ReloadCache() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}