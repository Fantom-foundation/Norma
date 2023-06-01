// Code generated by MockGen. DO NOT EDIT.
// Source: node.go

// Package driver is a generated GoMock package.
package driver

import (
	io "io"
	reflect "reflect"

	network "github.com/Fantom-foundation/Norma/driver/network"
	gomock "github.com/golang/mock/gomock"
)

// MockNode is a mock of Node interface.
type MockNode struct {
	ctrl     *gomock.Controller
	recorder *MockNodeMockRecorder
}

// MockNodeMockRecorder is the mock recorder for MockNode.
type MockNodeMockRecorder struct {
	mock *MockNode
}

// NewMockNode creates a new mock instance.
func NewMockNode(ctrl *gomock.Controller) *MockNode {
	mock := &MockNode{ctrl: ctrl}
	mock.recorder = &MockNodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNode) EXPECT() *MockNodeMockRecorder {
	return m.recorder
}

// Cleanup mocks base method.
func (m *MockNode) Cleanup() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cleanup")
	ret0, _ := ret[0].(error)
	return ret0
}

// Cleanup indicates an expected call of Cleanup.
func (mr *MockNodeMockRecorder) Cleanup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockNode)(nil).Cleanup))
}

// GetHttpServiceUrl mocks base method.
func (m *MockNode) GetHttpServiceUrl(arg0 *network.ServiceDescription) *URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHttpServiceUrl", arg0)
	ret0, _ := ret[0].(*URL)
	return ret0
}

// GetHttpServiceUrl indicates an expected call of GetHttpServiceUrl.
func (mr *MockNodeMockRecorder) GetHttpServiceUrl(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHttpServiceUrl", reflect.TypeOf((*MockNode)(nil).GetHttpServiceUrl), arg0)
}

// GetLabel mocks base method.
func (m *MockNode) GetLabel() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLabel")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetLabel indicates an expected call of GetLabel.
func (mr *MockNodeMockRecorder) GetLabel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLabel", reflect.TypeOf((*MockNode)(nil).GetLabel))
}

// GetNodeID mocks base method.
func (m *MockNode) GetNodeID() (NodeID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeID")
	ret0, _ := ret[0].(NodeID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeID indicates an expected call of GetNodeID.
func (mr *MockNodeMockRecorder) GetNodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeID", reflect.TypeOf((*MockNode)(nil).GetNodeID))
}

// GetWebsocketServiceUrl mocks base method.
func (m *MockNode) GetWebsocketServiceUrl(service *network.ServiceDescription) *URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWebsocketServiceUrl", service)
	ret0, _ := ret[0].(*URL)
	return ret0
}

// GetWebsocketServiceUrl indicates an expected call of GetWebsocketServiceUrl.
func (mr *MockNodeMockRecorder) GetWebsocketServiceUrl(service interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWebsocketServiceUrl", reflect.TypeOf((*MockNode)(nil).GetWebsocketServiceUrl), service)
}

// IsRunning mocks base method.
func (m *MockNode) IsRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRunning indicates an expected call of IsRunning.
func (mr *MockNodeMockRecorder) IsRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRunning", reflect.TypeOf((*MockNode)(nil).IsRunning))
}

// Stop mocks base method.
func (m *MockNode) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockNodeMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockNode)(nil).Stop))
}

// StreamLog mocks base method.
func (m *MockNode) StreamLog() (io.ReadCloser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StreamLog")
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StreamLog indicates an expected call of StreamLog.
func (mr *MockNodeMockRecorder) StreamLog() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StreamLog", reflect.TypeOf((*MockNode)(nil).StreamLog))
}
