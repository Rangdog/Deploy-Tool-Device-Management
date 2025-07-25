// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	entity "BE_Manage_device/internal/domain/entity"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// AssignmentRepository is an autogenerated mock type for the AssignmentRepository type
type AssignmentRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: assignment, tx
func (_m *AssignmentRepository) Create(assignment *entity.Assignments, tx *gorm.DB) (*entity.Assignments, error) {
	ret := _m.Called(assignment, tx)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *entity.Assignments
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.Assignments, *gorm.DB) (*entity.Assignments, error)); ok {
		return rf(assignment, tx)
	}
	if rf, ok := ret.Get(0).(func(*entity.Assignments, *gorm.DB) *entity.Assignments); ok {
		r0 = rf(assignment, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Assignments)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.Assignments, *gorm.DB) error); ok {
		r1 = rf(assignment, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAssignmentByAssetId provides a mock function with given fields: assetId
func (_m *AssignmentRepository) GetAssignmentByAssetId(assetId int64) (*entity.Assignments, error) {
	ret := _m.Called(assetId)

	if len(ret) == 0 {
		panic("no return value specified for GetAssignmentByAssetId")
	}

	var r0 *entity.Assignments
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*entity.Assignments, error)); ok {
		return rf(assetId)
	}
	if rf, ok := ret.Get(0).(func(int64) *entity.Assignments); ok {
		r0 = rf(assetId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Assignments)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(assetId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAssignmentById provides a mock function with given fields: id
func (_m *AssignmentRepository) GetAssignmentById(id int64) (*entity.Assignments, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetAssignmentById")
	}

	var r0 *entity.Assignments
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*entity.Assignments, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *entity.Assignments); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Assignments)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDB provides a mock function with no fields
func (_m *AssignmentRepository) GetDB() *gorm.DB {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDB")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Update provides a mock function with given fields: assignmentId, AssignBy, assetId, userId, departmentId, tx
func (_m *AssignmentRepository) Update(assignmentId int64, AssignBy int64, assetId int64, userId *int64, departmentId *int64, tx *gorm.DB) (*entity.Assignments, error) {
	ret := _m.Called(assignmentId, AssignBy, assetId, userId, departmentId, tx)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *entity.Assignments
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64, int64, *int64, *int64, *gorm.DB) (*entity.Assignments, error)); ok {
		return rf(assignmentId, AssignBy, assetId, userId, departmentId, tx)
	}
	if rf, ok := ret.Get(0).(func(int64, int64, int64, *int64, *int64, *gorm.DB) *entity.Assignments); ok {
		r0 = rf(assignmentId, AssignBy, assetId, userId, departmentId, tx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Assignments)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64, int64, *int64, *int64, *gorm.DB) error); ok {
		r1 = rf(assignmentId, AssignBy, assetId, userId, departmentId, tx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAssignmentRepository creates a new instance of AssignmentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAssignmentRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AssignmentRepository {
	mock := &AssignmentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
