// Code generated by mockery v2.32.2. DO NOT EDIT.

package mocks

import (
	api "github.com/DataDog/datadog-agent/pkg/eventmonitor/proto/api"
	mock "github.com/stretchr/testify/mock"
)

// EventMonitoringModuleServer is an autogenerated mock type for the EventMonitoringModuleServer type
type EventMonitoringModuleServer struct {
	mock.Mock
}

// GetProcessEvents provides a mock function with given fields: _a0, _a1
func (_m *EventMonitoringModuleServer) GetProcessEvents(_a0 *api.GetProcessEventParams, _a1 api.EventMonitoringModule_GetProcessEventsServer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*api.GetProcessEventParams, api.EventMonitoringModule_GetProcessEventsServer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mustEmbedUnimplementedEventMonitoringModuleServer provides a mock function with given fields:
func (_m *EventMonitoringModuleServer) mustEmbedUnimplementedEventMonitoringModuleServer() {
	_m.Called()
}

// NewEventMonitoringModuleServer creates a new instance of EventMonitoringModuleServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventMonitoringModuleServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *EventMonitoringModuleServer {
	mock := &EventMonitoringModuleServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
