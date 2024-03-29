// Code generated by MockGen. DO NOT EDIT.
// Source: src/repository/transfer/transfer.go

// Package mock_transfer is a generated GoMock package.
package mock_transfer

import (
	context "context"
	reflect "reflect"

	clientresponse "github.com/achwanyusuf/bricksvc/src/domain/clientresponse"
	entity "github.com/achwanyusuf/bricksvc/src/domain/entity"
	model "github.com/achwanyusuf/bricksvc/src/domain/model"
	gomock "github.com/golang/mock/gomock"
)

// MockTransferInterface is a mock of TransferInterface interface.
type MockTransferInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTransferInterfaceMockRecorder
}

// MockTransferInterfaceMockRecorder is the mock recorder for MockTransferInterface.
type MockTransferInterfaceMockRecorder struct {
	mock *MockTransferInterface
}

// NewMockTransferInterface creates a new mock instance.
func NewMockTransferInterface(ctrl *gomock.Controller) *MockTransferInterface {
	mock := &MockTransferInterface{ctrl: ctrl}
	mock.recorder = &MockTransferInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransferInterface) EXPECT() *MockTransferInterfaceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockTransferInterface) Delete(ctx context.Context, v *entity.TransferJob, id int64, isHardDelete bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, v, id, isHardDelete)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTransferInterfaceMockRecorder) Delete(ctx, v, id, isHardDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTransferInterface)(nil).Delete), ctx, v, id, isHardDelete)
}

// GetByParam mocks base method.
func (m *MockTransferInterface) GetByParam(ctx context.Context, cacheControl string, param *model.GetTransferJobsByParam) (entity.TransferJobSlice, model.Pagination, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByParam", ctx, cacheControl, param)
	ret0, _ := ret[0].(entity.TransferJobSlice)
	ret1, _ := ret[1].(model.Pagination)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByParam indicates an expected call of GetByParam.
func (mr *MockTransferInterfaceMockRecorder) GetByParam(ctx, cacheControl, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByParam", reflect.TypeOf((*MockTransferInterface)(nil).GetByParam), ctx, cacheControl, param)
}

// GetSingleByParam mocks base method.
func (m *MockTransferInterface) GetSingleByParam(ctx context.Context, cacheControl string, param *model.GetTransferJobByParam) (entity.TransferJob, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSingleByParam", ctx, cacheControl, param)
	ret0, _ := ret[0].(entity.TransferJob)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSingleByParam indicates an expected call of GetSingleByParam.
func (mr *MockTransferInterfaceMockRecorder) GetSingleByParam(ctx, cacheControl, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSingleByParam", reflect.TypeOf((*MockTransferInterface)(nil).GetSingleByParam), ctx, cacheControl, param)
}

// GetTransferByID mocks base method.
func (m *MockTransferInterface) GetTransferByID(id string) (clientresponse.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransferByID", id)
	ret0, _ := ret[0].(clientresponse.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransferByID indicates an expected call of GetTransferByID.
func (mr *MockTransferInterfaceMockRecorder) GetTransferByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransferByID", reflect.TypeOf((*MockTransferInterface)(nil).GetTransferByID), id)
}

// Insert mocks base method.
func (m *MockTransferInterface) Insert(ctx context.Context, data *entity.TransferJob) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockTransferInterfaceMockRecorder) Insert(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockTransferInterface)(nil).Insert), ctx, data)
}

// InsertTransfer mocks base method.
func (m *MockTransferInterface) InsertTransfer(ctx context.Context, data clientresponse.Transfer) (clientresponse.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertTransfer", ctx, data)
	ret0, _ := ret[0].(clientresponse.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertTransfer indicates an expected call of InsertTransfer.
func (mr *MockTransferInterfaceMockRecorder) InsertTransfer(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTransfer", reflect.TypeOf((*MockTransferInterface)(nil).InsertTransfer), ctx, data)
}

// Update mocks base method.
func (m *MockTransferInterface) Update(ctx context.Context, v *entity.TransferJob) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTransferInterfaceMockRecorder) Update(ctx, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTransferInterface)(nil).Update), ctx, v)
}
