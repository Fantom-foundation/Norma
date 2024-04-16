// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

// Code generated by MockGen. DO NOT EDIT.
// Source: host.go

// Package network is a generated GoMock package.
package network

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHost is a mock of Host interface.
type MockHost struct {
	ctrl     *gomock.Controller
	recorder *MockHostMockRecorder
}

// MockHostMockRecorder is the mock recorder for MockHost.
type MockHostMockRecorder struct {
	mock *MockHost
}

// NewMockHost creates a new mock instance.
func NewMockHost(ctrl *gomock.Controller) *MockHost {
	mock := &MockHost{ctrl: ctrl}
	mock.recorder = &MockHostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHost) EXPECT() *MockHostMockRecorder {
	return m.recorder
}

// Cleanup mocks base method.
func (m *MockHost) Cleanup() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cleanup")
	ret0, _ := ret[0].(error)
	return ret0
}

// Cleanup indicates an expected call of Cleanup.
func (mr *MockHostMockRecorder) Cleanup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockHost)(nil).Cleanup))
}

// GetAddressForService mocks base method.
func (m *MockHost) GetAddressForService(arg0 *ServiceDescription) *AddressPort {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressForService", arg0)
	ret0, _ := ret[0].(*AddressPort)
	return ret0
}

// GetAddressForService indicates an expected call of GetAddressForService.
func (mr *MockHostMockRecorder) GetAddressForService(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressForService", reflect.TypeOf((*MockHost)(nil).GetAddressForService), arg0)
}

// Hostname mocks base method.
func (m *MockHost) Hostname() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hostname")
	ret0, _ := ret[0].(string)
	return ret0
}

// Hostname indicates an expected call of Hostname.
func (mr *MockHostMockRecorder) Hostname() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hostname", reflect.TypeOf((*MockHost)(nil).Hostname))
}

// IsRunning mocks base method.
func (m *MockHost) IsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRunning indicates an expected call of IsRunning.
func (mr *MockHostMockRecorder) IsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockHost)(nil).IsRunning))
}

// SaveLogTo mocks base method.
func (m *MockHost) SaveLogTo(directory string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveLogTo", directory)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveLogTo indicates an expected call of SaveLogTo.
func (mr *MockHostMockRecorder) SaveLogTo(directory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveLogTo", reflect.TypeOf((*MockHost)(nil).SaveLogTo), directory)
}

// Stop mocks base method.
func (m *MockHost) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockHostMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockHost)(nil).Stop))
}

// StreamLog mocks base method.
func (m *MockHost) StreamLog() (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamLog")
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamLog indicates an expected call of StreamLog.
func (mr *MockHostMockRecorder) StreamLog() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamLog", reflect.TypeOf((*MockHost)(nil).StreamLog))
}
