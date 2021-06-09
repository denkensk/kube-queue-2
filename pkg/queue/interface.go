package queue

import (
	"github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
	"github.com/kube-queue/kube-queue/pkg/framework"
)

type MultiSchedulingQueue interface {
	Add(*v1alpha1.Queue) error
	Delete(*v1alpha1.Queue) error
	Update(*v1alpha1.Queue, *v1alpha1.Queue) error
	SortedQueue() []SchedulingQueue
	GetQueueByName(name string) (SchedulingQueue, bool)
}

type SchedulingQueue interface {
	Add(*v1alpha1.QueueUnit) error
	Delete(*v1alpha1.QueueUnit) error
	Update(*v1alpha1.QueueUnit, *v1alpha1.QueueUnit) error
	Pop() (*framework.QueueUnitInfo, error)
	TopUnit() (*framework.QueueUnitInfo, error)
	Name() string
	QueueInfo() *framework.QueueInfo
	Length() int
}
