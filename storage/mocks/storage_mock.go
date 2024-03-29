// Code generated by MockGen. DO NOT EDIT.
// Source: storage/storage.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	context "context"
	reflect "reflect"
	api "router-location-connecter/api"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddLocation mocks base method.
func (m *MockStorage) AddLocation(location *api.Location) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLocation", location)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLocation indicates an expected call of AddLocation.
func (mr *MockStorageMockRecorder) AddLocation(location interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLocation", reflect.TypeOf((*MockStorage)(nil).AddLocation), location)
}

// AddRouter mocks base method.
func (m *MockStorage) AddRouter(router *api.Router) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRouter", router)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRouter indicates an expected call of AddRouter.
func (mr *MockStorageMockRecorder) AddRouter(router interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRouter", reflect.TypeOf((*MockStorage)(nil).AddRouter), router)
}

// AddRouterLocationLink mocks base method.
func (m *MockStorage) AddRouterLocationLink(links *api.RouterLocationLink) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRouterLocationLink", links)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRouterLocationLink indicates an expected call of AddRouterLocationLink.
func (mr *MockStorageMockRecorder) AddRouterLocationLink(links interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRouterLocationLink", reflect.TypeOf((*MockStorage)(nil).AddRouterLocationLink), links)
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// FlushAll mocks base method.
func (m *MockStorage) FlushAll(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FlushAll", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// FlushAll indicates an expected call of FlushAll.
func (mr *MockStorageMockRecorder) FlushAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FlushAll", reflect.TypeOf((*MockStorage)(nil).FlushAll), ctx)
}

// GetLocation mocks base method.
func (m *MockStorage) GetLocation(id int) (*api.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocation", id)
	ret0, _ := ret[0].(*api.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocation indicates an expected call of GetLocation.
func (mr *MockStorageMockRecorder) GetLocation(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocation", reflect.TypeOf((*MockStorage)(nil).GetLocation), id)
}

// GetRouter mocks base method.
func (m *MockStorage) GetRouter(id int) (*api.Router, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRouter", id)
	ret0, _ := ret[0].(*api.Router)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRouter indicates an expected call of GetRouter.
func (mr *MockStorageMockRecorder) GetRouter(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRouter", reflect.TypeOf((*MockStorage)(nil).GetRouter), id)
}

// GetRouterLocationLink mocks base method.
func (m *MockStorage) GetRouterLocationLink(uniqueID string) (*api.RouterLocationLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRouterLocationLink", uniqueID)
	ret0, _ := ret[0].(*api.RouterLocationLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRouterLocationLink indicates an expected call of GetRouterLocationLink.
func (mr *MockStorageMockRecorder) GetRouterLocationLink(uniqueID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRouterLocationLink", reflect.TypeOf((*MockStorage)(nil).GetRouterLocationLink), uniqueID)
}
