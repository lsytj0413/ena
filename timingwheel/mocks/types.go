// Code generated by MockGen. DO NOT EDIT.
// Source: timingwheel/types.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	timingwheel "github.com/lsytj0413/ena/timingwheel"
)

// MockTimingWheel is a mock of TimingWheel interface.
type MockTimingWheel struct {
	ctrl     *gomock.Controller
	recorder *MockTimingWheelMockRecorder
}

// MockTimingWheelMockRecorder is the mock recorder for MockTimingWheel.
type MockTimingWheelMockRecorder struct {
	mock *MockTimingWheel
}

// NewMockTimingWheel creates a new mock instance.
func NewMockTimingWheel(ctrl *gomock.Controller) *MockTimingWheel {
	mock := &MockTimingWheel{ctrl: ctrl}
	mock.recorder = &MockTimingWheelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimingWheel) EXPECT() *MockTimingWheelMockRecorder {
	return m.recorder
}

// AfterFunc mocks base method.
func (m *MockTimingWheel) AfterFunc(d time.Duration, f timingwheel.Handler) (timingwheel.TimerTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AfterFunc", d, f)
	ret0, _ := ret[0].(timingwheel.TimerTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AfterFunc indicates an expected call of AfterFunc.
func (mr *MockTimingWheelMockRecorder) AfterFunc(d, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AfterFunc", reflect.TypeOf((*MockTimingWheel)(nil).AfterFunc), d, f)
}

// Start mocks base method.
func (m *MockTimingWheel) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockTimingWheelMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockTimingWheel)(nil).Start))
}

// Stop mocks base method.
func (m *MockTimingWheel) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockTimingWheelMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTimingWheel)(nil).Stop))
}

// TickFunc mocks base method.
func (m *MockTimingWheel) TickFunc(d time.Duration, f timingwheel.Handler) (timingwheel.TimerTask, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TickFunc", d, f)
	ret0, _ := ret[0].(timingwheel.TimerTask)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TickFunc indicates an expected call of TickFunc.
func (mr *MockTimingWheelMockRecorder) TickFunc(d, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TickFunc", reflect.TypeOf((*MockTimingWheel)(nil).TickFunc), d, f)
}

// MockTimerTask is a mock of TimerTask interface.
type MockTimerTask struct {
	ctrl     *gomock.Controller
	recorder *MockTimerTaskMockRecorder
}

// MockTimerTaskMockRecorder is the mock recorder for MockTimerTask.
type MockTimerTaskMockRecorder struct {
	mock *MockTimerTask
}

// NewMockTimerTask creates a new mock instance.
func NewMockTimerTask(ctrl *gomock.Controller) *MockTimerTask {
	mock := &MockTimerTask{ctrl: ctrl}
	mock.recorder = &MockTimerTaskMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimerTask) EXPECT() *MockTimerTaskMockRecorder {
	return m.recorder
}

// Stop mocks base method.
func (m *MockTimerTask) Stop() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockTimerTaskMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTimerTask)(nil).Stop))
}
