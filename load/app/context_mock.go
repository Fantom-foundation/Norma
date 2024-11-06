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
// Source: context.go
//
// Generated by this command:
//
//	mockgen -source context.go -destination context_mock.go -package app
//

// Package app is a generated GoMock package.
package app

import (
	big "math/big"
	reflect "reflect"

	rpc "github.com/Fantom-foundation/Norma/driver/rpc"
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	gomock "go.uber.org/mock/gomock"
)

// MockAppContext is a mock of AppContext interface.
type MockAppContext struct {
	ctrl     *gomock.Controller
	recorder *MockAppContextMockRecorder
}

// MockAppContextMockRecorder is the mock recorder for MockAppContext.
type MockAppContextMockRecorder struct {
	mock *MockAppContext
}

// NewMockAppContext creates a new mock instance.
func NewMockAppContext(ctrl *gomock.Controller) *MockAppContext {
	mock := &MockAppContext{ctrl: ctrl}
	mock.recorder = &MockAppContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppContext) EXPECT() *MockAppContextMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockAppContext) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockAppContextMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAppContext)(nil).Close))
}

// FundAccounts mocks base method.
func (m *MockAppContext) FundAccounts(accounts []common.Address, value *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FundAccounts", accounts, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// FundAccounts indicates an expected call of FundAccounts.
func (mr *MockAppContextMockRecorder) FundAccounts(accounts, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FundAccounts", reflect.TypeOf((*MockAppContext)(nil).FundAccounts), accounts, value)
}

// GetClient mocks base method.
func (m *MockAppContext) GetClient() rpc.RpcClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClient")
	ret0, _ := ret[0].(rpc.RpcClient)
	return ret0
}

// GetClient indicates an expected call of GetClient.
func (mr *MockAppContextMockRecorder) GetClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockAppContext)(nil).GetClient))
}

// GetReceipt mocks base method.
func (m *MockAppContext) GetReceipt(txHash common.Hash) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReceipt", txHash)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReceipt indicates an expected call of GetReceipt.
func (mr *MockAppContextMockRecorder) GetReceipt(txHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReceipt", reflect.TypeOf((*MockAppContext)(nil).GetReceipt), txHash)
}

// GetTransactOptions mocks base method.
func (m *MockAppContext) GetTransactOptions(account *Account) (*bind.TransactOpts, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactOptions", account)
	ret0, _ := ret[0].(*bind.TransactOpts)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactOptions indicates an expected call of GetTransactOptions.
func (mr *MockAppContextMockRecorder) GetTransactOptions(account any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactOptions", reflect.TypeOf((*MockAppContext)(nil).GetTransactOptions), account)
}

// GetTreasure mocks base method.
func (m *MockAppContext) GetTreasure() *Account {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTreasure")
	ret0, _ := ret[0].(*Account)
	return ret0
}

// GetTreasure indicates an expected call of GetTreasure.
func (mr *MockAppContextMockRecorder) GetTreasure() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTreasure", reflect.TypeOf((*MockAppContext)(nil).GetTreasure))
}

// Run mocks base method.
func (m *MockAppContext) Run(operation func(*bind.TransactOpts) (*types.Transaction, error)) (*types.Receipt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", operation)
	ret0, _ := ret[0].(*types.Receipt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run.
func (mr *MockAppContextMockRecorder) Run(operation any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockAppContext)(nil).Run), operation)
}
