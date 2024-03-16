// Code generated by MockGen. DO NOT EDIT.
// Source: src/usecase/bank/bank.go

// Package mock_bank is a generated GoMock package.
package mock_bank

import (
	context "context"
	reflect "reflect"

	clientresponse "github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	model "github.com/achwanyusuf/bricksvc/src/domain/model"
	gomock "github.com/golang/mock/gomock"
)

// MockBankInterface is a mock of BankInterface interface.
type MockBankInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBankInterfaceMockRecorder
}

// MockBankInterfaceMockRecorder is the mock recorder for MockBankInterface.
type MockBankInterfaceMockRecorder struct {
	mock *MockBankInterface
}

// NewMockBankInterface creates a new mock instance.
func NewMockBankInterface(ctrl *gomock.Controller) *MockBankInterface {
	mock := &MockBankInterface{ctrl: ctrl}
	mock.recorder = &MockBankInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBankInterface) EXPECT() *MockBankInterfaceMockRecorder {
	return m.recorder
}

// GetBankAccount mocks base method.
func (m *MockBankInterface) GetBankAccount(ctx context.Context, cacheControl string, v model.GetBankAccount, apikey string) (clientresponse.BankAccount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBankAccount", ctx, cacheControl, v, apikey)
	ret0, _ := ret[0].(clientresponse.BankAccount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBankAccount indicates an expected call of GetBankAccount.
func (mr *MockBankInterfaceMockRecorder) GetBankAccount(ctx, cacheControl, v, apikey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBankAccount", reflect.TypeOf((*MockBankInterface)(nil).GetBankAccount), ctx, cacheControl, v, apikey)
}
