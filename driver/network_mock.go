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
// Source: network.go
//
// Generated by this command:
//
//	mockgen -source network.go -destination network_mock.go -package driver
//

// Package driver is a generated GoMock package.
package driver

import (
	reflect "reflect"

	rpc "github.com/Fantom-foundation/Norma/driver/rpc"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockNetwork is a mock of Network interface.
type MockNetwork struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkMockRecorder
}

// MockNetworkMockRecorder is the mock recorder for MockNetwork.
type MockNetworkMockRecorder struct {
	mock *MockNetwork
}

// NewMockNetwork creates a new mock instance.
func NewMockNetwork(ctrl *gomock.Controller) *MockNetwork {
	mock := &MockNetwork{ctrl: ctrl}
	mock.recorder = &MockNetworkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetwork) EXPECT() *MockNetworkMockRecorder {
	return m.recorder
}

// CreateApplication mocks base method.
func (m *MockNetwork) CreateApplication(config *ApplicationConfig) (Application, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApplication", config)
	ret0, _ := ret[0].(Application)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateApplication indicates an expected call of CreateApplication.
func (mr *MockNetworkMockRecorder) CreateApplication(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApplication", reflect.TypeOf((*MockNetwork)(nil).CreateApplication), config)
}

// CreateNode mocks base method.
func (m *MockNetwork) CreateNode(config *NodeConfig) (Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNode", config)
	ret0, _ := ret[0].(Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNode indicates an expected call of CreateNode.
func (mr *MockNetworkMockRecorder) CreateNode(config any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNode", reflect.TypeOf((*MockNetwork)(nil).CreateNode), config)
}

// DialRandomRpc mocks base method.
func (m *MockNetwork) DialRandomRpc() (rpc.RpcClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DialRandomRpc")
	ret0, _ := ret[0].(rpc.RpcClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DialRandomRpc indicates an expected call of DialRandomRpc.
func (mr *MockNetworkMockRecorder) DialRandomRpc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DialRandomRpc", reflect.TypeOf((*MockNetwork)(nil).DialRandomRpc))
}

// GetActiveApplications mocks base method.
func (m *MockNetwork) GetActiveApplications() []Application {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveApplications")
	ret0, _ := ret[0].([]Application)
	return ret0
}

// GetActiveApplications indicates an expected call of GetActiveApplications.
func (mr *MockNetworkMockRecorder) GetActiveApplications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveApplications", reflect.TypeOf((*MockNetwork)(nil).GetActiveApplications))
}

// GetActiveNodes mocks base method.
func (m *MockNetwork) GetActiveNodes() []Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveNodes")
	ret0, _ := ret[0].([]Node)
	return ret0
}

// GetActiveNodes indicates an expected call of GetActiveNodes.
func (mr *MockNetworkMockRecorder) GetActiveNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveNodes", reflect.TypeOf((*MockNetwork)(nil).GetActiveNodes))
}

// RegisterListener mocks base method.
func (m *MockNetwork) RegisterListener(arg0 NetworkListener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterListener", arg0)
}

// RegisterListener indicates an expected call of RegisterListener.
func (mr *MockNetworkMockRecorder) RegisterListener(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterListener", reflect.TypeOf((*MockNetwork)(nil).RegisterListener), arg0)
}

// RemoveNode mocks base method.
func (m *MockNetwork) RemoveNode(arg0 Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNode indicates an expected call of RemoveNode.
func (mr *MockNetworkMockRecorder) RemoveNode(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNode", reflect.TypeOf((*MockNetwork)(nil).RemoveNode), arg0)
}

// SendTransaction mocks base method.
func (m *MockNetwork) SendTransaction(tx *types.Transaction) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendTransaction", tx)
}

// SendTransaction indicates an expected call of SendTransaction.
func (mr *MockNetworkMockRecorder) SendTransaction(tx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTransaction", reflect.TypeOf((*MockNetwork)(nil).SendTransaction), tx)
}

// Shutdown mocks base method.
func (m *MockNetwork) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockNetworkMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockNetwork)(nil).Shutdown))
}

// UnregisterListener mocks base method.
func (m *MockNetwork) UnregisterListener(arg0 NetworkListener) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UnregisterListener", arg0)
}

// UnregisterListener indicates an expected call of UnregisterListener.
func (mr *MockNetworkMockRecorder) UnregisterListener(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterListener", reflect.TypeOf((*MockNetwork)(nil).UnregisterListener), arg0)
}

// MockNetworkListener is a mock of NetworkListener interface.
type MockNetworkListener struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkListenerMockRecorder
}

// MockNetworkListenerMockRecorder is the mock recorder for MockNetworkListener.
type MockNetworkListenerMockRecorder struct {
	mock *MockNetworkListener
}

// NewMockNetworkListener creates a new mock instance.
func NewMockNetworkListener(ctrl *gomock.Controller) *MockNetworkListener {
	mock := &MockNetworkListener{ctrl: ctrl}
	mock.recorder = &MockNetworkListenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkListener) EXPECT() *MockNetworkListenerMockRecorder {
	return m.recorder
}

// AfterApplicationCreation mocks base method.
func (m *MockNetworkListener) AfterApplicationCreation(arg0 Application) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AfterApplicationCreation", arg0)
}

// AfterApplicationCreation indicates an expected call of AfterApplicationCreation.
func (mr *MockNetworkListenerMockRecorder) AfterApplicationCreation(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterApplicationCreation", reflect.TypeOf((*MockNetworkListener)(nil).AfterApplicationCreation), arg0)
}

// AfterNodeCreation mocks base method.
func (m *MockNetworkListener) AfterNodeCreation(arg0 Node) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AfterNodeCreation", arg0)
}

// AfterNodeCreation indicates an expected call of AfterNodeCreation.
func (mr *MockNetworkListenerMockRecorder) AfterNodeCreation(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterNodeCreation", reflect.TypeOf((*MockNetworkListener)(nil).AfterNodeCreation), arg0)
}

// AfterNodeRemoval mocks base method.
func (m *MockNetworkListener) AfterNodeRemoval(arg0 Node) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AfterNodeRemoval", arg0)
}

// AfterNodeRemoval indicates an expected call of AfterNodeRemoval.
func (mr *MockNetworkListenerMockRecorder) AfterNodeRemoval(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterNodeRemoval", reflect.TypeOf((*MockNetworkListener)(nil).AfterNodeRemoval), arg0)
}
