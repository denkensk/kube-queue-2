package framework

import (
	"context"

	"k8s.io/api/core/v1"
)

// Code is the Status code/type which is returned from plugins.
type Code int

const (
	// Success means that plugin ran correctly and found pod schedulable.
	// NOTE: A nil status is also considered as "Success".
	Success Code = iota
	// Error is used for internal plugin errors, unexpected input, etc.
	Error
	// Unschedulable is used when a plugin finds a pod unschedulable. The scheduler might attempt to
	// preempt other pods to get this pod scheduled. Use UnschedulableAndUnresolvable to make the
	// scheduler skip preemption.
	// The accompanying status message should explain why the pod is unschedulable.
	Unschedulable
	// UnschedulableAndUnresolvable is used when a PreFilter plugin finds a pod unschedulable and
	// preemption would not change anything. Plugins should return Unschedulable if it is possible
	// that the pod can get scheduled with preemption.
	// The accompanying status message should explain why the pod is unschedulable.
	UnschedulableAndUnresolvable
	// Wait is used when a Permit plugin finds a pod scheduling should wait.
	Wait
	// Skip is used when a Bind plugin chooses to skip binding.
	Skip
)

type Framework interface {
	// QueueSortFunc returns the function to sort pods in scheduling queue
	MultiQueueSortFunc() MultiQueueLessFunc
	QueueSortFuncMap() map[string]QueueLessFunc
	RunFilterPlugins(*QueueUnitInfo) *Status
	RunScorePlugins() (int64, bool)
}

type Status struct {
	message string
	code    Code
}

// NewStatus makes a Status out of the given arguments and returns its pointer.
func NewStatus(code Code, message string) *Status {
	s := &Status{
		code:    code,
		message: message,
	}
	return s
}

func (s *Status) Code() Code {
	return s.code
}

// Plugin is the parent type for all the scheduling framework plugins.
type Plugin interface {
	Name() string
}

type MultiQueueSortPlugin interface {
	Plugin
	MultiQueueLess(*QueueInfo, *QueueInfo) bool
}

type MultiQueueLessFunc func(*QueueInfo, *QueueInfo) bool

type QueueSortPlugin interface {
	Plugin
	QueueLess(*QueueUnitInfo, *QueueUnitInfo) bool
}

type QueueLessFunc func(*QueueUnitInfo, *QueueUnitInfo) bool

type FilterPlugin interface {
	Plugin

	Filter(ctx context.Context, QueueUnit *QueueUnitInfo) *Status
}

// ScorePlugin is an interface that must be implemented by "Score" plugins to rank
// nodes that passed the filtering phase.
type ScorePlugin interface {
	Plugin
	Score(ctx context.Context, p *v1.Pod) (int64, bool)
}

type Handle interface {
}
