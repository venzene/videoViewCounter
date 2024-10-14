// Code generated by MockGen. DO NOT EDIT.
// Source: repository/viewrepository/repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	model "view_count/model"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetAllViews mocks base method.
func (m *MockRepository) GetAllViews(ctx context.Context) ([]model.VideoInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllViews", ctx)
	ret0, _ := ret[0].([]model.VideoInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllViews indicates an expected call of GetAllViews.
func (mr *MockRepositoryMockRecorder) GetAllViews(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllViews", reflect.TypeOf((*MockRepository)(nil).GetAllViews), ctx)
}

// GetRecentVideos mocks base method.
func (m *MockRepository) GetRecentVideos(ctx context.Context, n int) ([]model.VideoInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecentVideos", ctx, n)
	ret0, _ := ret[0].([]model.VideoInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecentVideos indicates an expected call of GetRecentVideos.
func (mr *MockRepositoryMockRecorder) GetRecentVideos(ctx, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecentVideos", reflect.TypeOf((*MockRepository)(nil).GetRecentVideos), ctx, n)
}

// GetTopVideos mocks base method.
func (m *MockRepository) GetTopVideos(ctx context.Context, n int) ([]model.VideoInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopVideos", ctx, n)
	ret0, _ := ret[0].([]model.VideoInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopVideos indicates an expected call of GetTopVideos.
func (mr *MockRepositoryMockRecorder) GetTopVideos(ctx, n interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopVideos", reflect.TypeOf((*MockRepository)(nil).GetTopVideos), ctx, n)
}

// GetView mocks base method.
func (m *MockRepository) GetView(ctx context.Context, videoId string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetView", ctx, videoId)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetView indicates an expected call of GetView.
func (mr *MockRepositoryMockRecorder) GetView(ctx, videoId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetView", reflect.TypeOf((*MockRepository)(nil).GetView), ctx, videoId)
}

// Increment mocks base method.
func (m *MockRepository) Increment(ctx context.Context, videoId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Increment", ctx, videoId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Increment indicates an expected call of Increment.
func (mr *MockRepositoryMockRecorder) Increment(ctx, videoId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockRepository)(nil).Increment), ctx, videoId)
}
