// Code generated by mockery v1.0.0. DO NOT EDIT.

package clusterfeatureadapter

import context "context"
import mock "github.com/stretchr/testify/mock"

// MockClusterGetter is an autogenerated mock type for the ClusterGetter type
type MockClusterGetter struct {
	mock.Mock
}

// GetClusterByIDOnly provides a mock function with given fields: ctx, clusterID
func (_m *MockClusterGetter) GetClusterByIDOnly(ctx context.Context, clusterID uint) (Cluster, error) {
	ret := _m.Called(ctx, clusterID)

	var r0 Cluster
	if rf, ok := ret.Get(0).(func(context.Context, uint) Cluster); ok {
		r0 = rf(ctx, clusterID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Cluster)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, clusterID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
