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
// Source: node_log_provider.go

// Package monitoring is a generated GoMock package.
package monitoring

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogListener is a mock of LogListener interface.
type MockLogListener struct {
	ctrl     *gomock.Controller
	recorder *MockLogListenerMockRecorder
}

// MockLogListenerMockRecorder is the mock recorder for MockLogListener.
type MockLogListenerMockRecorder struct {
	mock *MockLogListener
}

// NewMockLogListener creates a new mock instance.
func NewMockLogListener(ctrl *gomock.Controller) *MockLogListener {
	mock := &MockLogListener{ctrl: ctrl}
	mock.recorder = &MockLogListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogListener) EXPECT() *MockLogListenerMockRecorder {
	return m.recorder
}

// OnBlock mocks base method.
func (m *MockLogListener) OnBlock(node Node, block Block) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnBlock", node, block)
}

// OnBlock indicates an expected call of OnBlock.
func (mr *MockLogListenerMockRecorder) OnBlock(node, block interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnBlock", reflect.TypeOf((*MockLogListener)(nil).OnBlock), node, block)
}

// MockNodeLogProvider is a mock of NodeLogProvider interface.
type MockNodeLogProvider struct {
	ctrl     *gomock.Controller
	recorder *MockNodeLogProviderMockRecorder
}

// MockNodeLogProviderMockRecorder is the mock recorder for MockNodeLogProvider.
type MockNodeLogProviderMockRecorder struct {
	mock *MockNodeLogProvider
}

// NewMockNodeLogProvider creates a new mock instance.
func NewMockNodeLogProvider(ctrl *gomock.Controller) *MockNodeLogProvider {
	mock := &MockNodeLogProvider{ctrl: ctrl}
	mock.recorder = &MockNodeLogProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeLogProvider) EXPECT() *MockNodeLogProviderMockRecorder {
	return m.recorder
}

// RegisterLogListener mocks base method.
func (m *MockNodeLogProvider) RegisterLogListener(listener LogListener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterLogListener", listener)
}

// RegisterLogListener indicates an expected call of RegisterLogListener.
func (mr *MockNodeLogProviderMockRecorder) RegisterLogListener(listener interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterLogListener", reflect.TypeOf((*MockNodeLogProvider)(nil).RegisterLogListener), listener)
}

// UnregisterLogListener mocks base method.
func (m *MockNodeLogProvider) UnregisterLogListener(listener LogListener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnregisterLogListener", listener)
}

// UnregisterLogListener indicates an expected call of UnregisterLogListener.
func (mr *MockNodeLogProviderMockRecorder) UnregisterLogListener(listener interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterLogListener", reflect.TypeOf((*MockNodeLogProvider)(nil).UnregisterLogListener), listener)
}
