package framework

import (
	"github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
)

type QueueInfo struct {
	name     string
	priority int32
	queue    *v1alpha1.Queue
}

func (q *QueueInfo) Name() string {
	return q.name
}

func (q *QueueInfo) Queue() *v1alpha1.Queue {
	return q.queue
}

func (q *QueueInfo) Priority() int32 {
	return q.queue.Spec.Priority
}

type QueueUnitInfo struct {
	name string
	unit *v1alpha1.QueueUnit
}

func (u *QueueUnitInfo) Name() string {
	return u.name
}

func (u *QueueUnitInfo) Unit() *v1alpha1.QueueUnit {
	return u.unit
}

func NewQueueUnitInfo(unit *v1alpha1.QueueUnit) *QueueUnitInfo {
	return &QueueUnitInfo{
		name: unit.Name,
		unit: unit,
	}
}
