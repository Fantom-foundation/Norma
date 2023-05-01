// Code generated by MockGen. DO NOT EDIT.
// Source: node.go

// Package driver is a generated GoMock package.
package driver

import (
	reflect "reflect"

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

// GetRpcServiceUrl mocks base method.
func (m *MockNode) GetRpcServiceUrl() *URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRpcServiceUrl")
	ret0, _ := ret[0].(*URL)
	return ret0
}

// GetRpcServiceUrl indicates an expected call of GetRpcServiceUrl.
func (mr *MockNodeMockRecorder) GetRpcServiceUrl() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRpcServiceUrl", reflect.TypeOf((*MockNode)(nil).GetRpcServiceUrl))
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
