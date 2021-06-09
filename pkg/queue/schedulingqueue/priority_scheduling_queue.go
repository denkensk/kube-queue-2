package schedulingqueue

import (
	"sync"

	"github.com/kube-queue/kube-queue/pkg/queue"

	"github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
	"github.com/kube-queue/kube-queue/pkg/framework"
	"github.com/kube-queue/kube-queue/pkg/queue/heap"
)

type PrioritySchedulingQueue struct {
	name       string
	pluginName string
	fw         framework.Framework
	items      *heap.Heap
	lock       sync.RWMutex
	queue      *framework.QueueInfo
}

func NewPrioritySchedulingQueue(fw framework.Framework, name string, pluginName string) queue.SchedulingQueue {
	queueSortFuncMap := fw.QueueSortFuncMap()
	lessFn := queueSortFuncMap[pluginName]
	comp := func(queueUnitInfo1, queueUnitInfo2 interface{}) bool {
		quInfo1 := queueUnitInfo1.(*framework.QueueUnitInfo)
		quInfo2 := queueUnitInfo2.(*framework.QueueUnitInfo)
		return lessFn(quInfo1, quInfo2)
	}

	q := &PrioritySchedulingQueue{
		fw:         fw,
		name:       name,
		pluginName: pluginName,
		items:      heap.New(unitInfoKeyFunc, comp),
	}

	return q
}

func (p *PrioritySchedulingQueue) Add(q *v1alpha1.QueueUnit) error {
	info := framework.NewQueueUnitInfo(q)
	err := p.items.Add(info)
	return err
}

func (p *PrioritySchedulingQueue) Delete(q *v1alpha1.QueueUnit) error {
	info := framework.NewQueueUnitInfo(q)
	return p.items.Delete(info)
}

func (p *PrioritySchedulingQueue) Update(old *v1alpha1.QueueUnit, new *v1alpha1.QueueUnit) error {
	info := framework.NewQueueUnitInfo(new)
	return p.items.Update(info)
}

func (p *PrioritySchedulingQueue) Pop() (*framework.QueueUnitInfo, error) {
	obj, err := p.items.Pop()
	u := obj.(*framework.QueueUnitInfo)
	return u, err
}

func (p *PrioritySchedulingQueue) TopUnit() (*framework.QueueUnitInfo, error) {
	return p.Pop()
}

func (p *PrioritySchedulingQueue) Name() string {
	return p.name
}

func (p *PrioritySchedulingQueue) QueueInfo() *framework.QueueInfo {
	return p.queue
}

func (p *PrioritySchedulingQueue) Length() int {
	return p.items.Len()
}

func unitInfoKeyFunc(obj interface{}) (string, error) {
	unitInfo := obj.(*framework.QueueUnitInfo)
	return unitInfo.Name(), nil
}
