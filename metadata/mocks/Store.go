package mocks

import "github.com/mistifyio/mistify-image-service/metadata"
import "github.com/stretchr/testify/mock"

// Store mocked by mockery
type Store struct {
	mock.Mock
}

// Init mocked by mockery
func (_m *Store) Init(_a0 []byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown mocked by mockery
func (_m *Store) Shutdown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List mocked by mockery
func (_m *Store) List(_a0 string) ([]*metadata.Image, error) {
	ret := _m.Called(_a0)

	var r0 []*metadata.Image
	if rf, ok := ret.Get(0).(func(string) []*metadata.Image); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*metadata.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID mocked by mockery
func (_m *Store) GetByID(_a0 string) (*metadata.Image, error) {
	ret := _m.Called(_a0)

	var r0 *metadata.Image
	if rf, ok := ret.Get(0).(func(string) *metadata.Image); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBySource mocked by mockery
func (_m *Store) GetBySource(_a0 string) (*metadata.Image, error) {
	ret := _m.Called(_a0)

	var r0 *metadata.Image
	if rf, ok := ret.Get(0).(func(string) *metadata.Image); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*metadata.Image)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Put mocked by mockery
func (_m *Store) Put(_a0 *metadata.Image) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*metadata.Image) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete mocked by mockery
func (_m *Store) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
