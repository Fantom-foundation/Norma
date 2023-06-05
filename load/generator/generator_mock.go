// Code generated by MockGen. DO NOT EDIT.
// Source: generator.go

// Package generator is a generated GoMock package.
package generator

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTransactionGenerator is a mock of TransactionGenerator interface.
type MockTransactionGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionGeneratorMockRecorder
}

// MockTransactionGeneratorMockRecorder is the mock recorder for MockTransactionGenerator.
type MockTransactionGeneratorMockRecorder struct {
	mock *MockTransactionGenerator
}

// NewMockTransactionGenerator creates a new mock instance.
func NewMockTransactionGenerator(ctrl *gomock.Controller) *MockTransactionGenerator {
	mock := &MockTransactionGenerator{ctrl: ctrl}
	mock.recorder = &MockTransactionGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionGenerator) EXPECT() *MockTransactionGeneratorMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockTransactionGenerator) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockTransactionGeneratorMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTransactionGenerator)(nil).Close))
}

// SendTx mocks base method.
func (m *MockTransactionGenerator) SendTx() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendTx")
	ret0, _ := ret[0].(error)
	return ret0
}

// SendTx indicates an expected call of SendTx.
func (mr *MockTransactionGeneratorMockRecorder) SendTx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTx", reflect.TypeOf((*MockTransactionGenerator)(nil).SendTx))
}

// MockTransactionGeneratorFactory is a mock of TransactionGeneratorFactory interface.
type MockTransactionGeneratorFactory struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionGeneratorFactoryMockRecorder
}

// MockTransactionGeneratorFactoryMockRecorder is the mock recorder for MockTransactionGeneratorFactory.
type MockTransactionGeneratorFactoryMockRecorder struct {
	mock *MockTransactionGeneratorFactory
}

// NewMockTransactionGeneratorFactory creates a new mock instance.
func NewMockTransactionGeneratorFactory(ctrl *gomock.Controller) *MockTransactionGeneratorFactory {
	mock := &MockTransactionGeneratorFactory{ctrl: ctrl}
	mock.recorder = &MockTransactionGeneratorFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionGeneratorFactory) EXPECT() *MockTransactionGeneratorFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTransactionGeneratorFactory) Create() (TransactionGenerator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(TransactionGenerator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionGeneratorFactoryMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactionGeneratorFactory)(nil).Create))
}

// MockTransactionGeneratorFactoryWithStats is a mock of TransactionGeneratorFactoryWithStats interface.
type MockTransactionGeneratorFactoryWithStats struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionGeneratorFactoryWithStatsMockRecorder
}

// MockTransactionGeneratorFactoryWithStatsMockRecorder is the mock recorder for MockTransactionGeneratorFactoryWithStats.
type MockTransactionGeneratorFactoryWithStatsMockRecorder struct {
	mock *MockTransactionGeneratorFactoryWithStats
}

// NewMockTransactionGeneratorFactoryWithStats creates a new mock instance.
func NewMockTransactionGeneratorFactoryWithStats(ctrl *gomock.Controller) *MockTransactionGeneratorFactoryWithStats {
	mock := &MockTransactionGeneratorFactoryWithStats{ctrl: ctrl}
	mock.recorder = &MockTransactionGeneratorFactoryWithStatsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionGeneratorFactoryWithStats) EXPECT() *MockTransactionGeneratorFactoryWithStatsMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTransactionGeneratorFactoryWithStats) Create() (TransactionGenerator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create")
	ret0, _ := ret[0].(TransactionGenerator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTransactionGeneratorFactoryWithStatsMockRecorder) Create() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTransactionGeneratorFactoryWithStats)(nil).Create))
}

// GetAmountOfReceivedTxs mocks base method.
func (m *MockTransactionGeneratorFactoryWithStats) GetAmountOfReceivedTxs() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountOfReceivedTxs")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAmountOfReceivedTxs indicates an expected call of GetAmountOfReceivedTxs.
func (mr *MockTransactionGeneratorFactoryWithStatsMockRecorder) GetAmountOfReceivedTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountOfReceivedTxs", reflect.TypeOf((*MockTransactionGeneratorFactoryWithStats)(nil).GetAmountOfReceivedTxs))
}

// GetAmountOfSentTxs mocks base method.
func (m *MockTransactionGeneratorFactoryWithStats) GetAmountOfSentTxs() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountOfSentTxs")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetAmountOfSentTxs indicates an expected call of GetAmountOfSentTxs.
func (mr *MockTransactionGeneratorFactoryWithStatsMockRecorder) GetAmountOfSentTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountOfSentTxs", reflect.TypeOf((*MockTransactionGeneratorFactoryWithStats)(nil).GetAmountOfSentTxs))
}

// MockTransactionCounts is a mock of TransactionCounts interface.
type MockTransactionCounts struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionCountsMockRecorder
}

// MockTransactionCountsMockRecorder is the mock recorder for MockTransactionCounts.
type MockTransactionCountsMockRecorder struct {
	mock *MockTransactionCounts
}

// NewMockTransactionCounts creates a new mock instance.
func NewMockTransactionCounts(ctrl *gomock.Controller) *MockTransactionCounts {
	mock := &MockTransactionCounts{ctrl: ctrl}
	mock.recorder = &MockTransactionCountsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionCounts) EXPECT() *MockTransactionCountsMockRecorder {
	return m.recorder
}

// GetAmountOfReceivedTxs mocks base method.
func (m *MockTransactionCounts) GetAmountOfReceivedTxs() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountOfReceivedTxs")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAmountOfReceivedTxs indicates an expected call of GetAmountOfReceivedTxs.
func (mr *MockTransactionCountsMockRecorder) GetAmountOfReceivedTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountOfReceivedTxs", reflect.TypeOf((*MockTransactionCounts)(nil).GetAmountOfReceivedTxs))
}

// GetAmountOfSentTxs mocks base method.
func (m *MockTransactionCounts) GetAmountOfSentTxs() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountOfSentTxs")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetAmountOfSentTxs indicates an expected call of GetAmountOfSentTxs.
func (mr *MockTransactionCountsMockRecorder) GetAmountOfSentTxs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountOfSentTxs", reflect.TypeOf((*MockTransactionCounts)(nil).GetAmountOfSentTxs))
}
