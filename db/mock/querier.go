// Code generated by MockGen. DO NOT EDIT.
// Source: db/sqlc/querier.go

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	sqlc "promova-test-task/db/sqlc"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreatePost mocks base method.
func (m *MockQuerier) CreatePost(ctx context.Context, arg sqlc.CreatePostParams) (sqlc.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, arg)
	ret0, _ := ret[0].(sqlc.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockQuerierMockRecorder) CreatePost(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockQuerier)(nil).CreatePost), ctx, arg)
}

// DeletePost mocks base method.
func (m *MockQuerier) DeletePost(ctx context.Context, id int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockQuerierMockRecorder) DeletePost(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockQuerier)(nil).DeletePost), ctx, id)
}

// GetPostById mocks base method.
func (m *MockQuerier) GetPostById(ctx context.Context, id int32) (sqlc.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostById", ctx, id)
	ret0, _ := ret[0].(sqlc.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostById indicates an expected call of GetPostById.
func (mr *MockQuerierMockRecorder) GetPostById(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostById", reflect.TypeOf((*MockQuerier)(nil).GetPostById), ctx, id)
}

// GetPosts mocks base method.
func (m *MockQuerier) GetPosts(ctx context.Context) ([]sqlc.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPosts", ctx)
	ret0, _ := ret[0].([]sqlc.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPosts indicates an expected call of GetPosts.
func (mr *MockQuerierMockRecorder) GetPosts(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPosts", reflect.TypeOf((*MockQuerier)(nil).GetPosts), ctx)
}

// UpdatePostById mocks base method.
func (m *MockQuerier) UpdatePostById(ctx context.Context, arg sqlc.UpdatePostByIdParams) (sqlc.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePostById", ctx, arg)
	ret0, _ := ret[0].(sqlc.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePostById indicates an expected call of UpdatePostById.
func (mr *MockQuerierMockRecorder) UpdatePostById(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePostById", reflect.TypeOf((*MockQuerier)(nil).UpdatePostById), ctx, arg)
}
