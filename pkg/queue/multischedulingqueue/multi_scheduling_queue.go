package multischedulingqueue

import (
	"sort"
	"sync"

	"github.com/kube-queue/kube-queue/pkg/queue/schedulingqueue"

	"github.com/kube-queue/kube-queue/pkg/queue"

	"github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
	"github.com/kube-queue/kube-queue/pkg/framework"
)

type multiSchedulingQueue struct {
	sync.RWMutex
	fw       framework.Framework
	queueMap map[string]queue.SchedulingQueue
	lessFunc framework.MultiQueueLessFunc
}

func NewMultiSchedulingQueue(fw framework.Framework) (queue.MultiSchedulingQueue, error) {

	mq := &multiSchedulingQueue{
		fw:       fw,
		queueMap: make(map[string]queue.SchedulingQueue),
		lessFunc: fw.MultiQueueSortFunc(),
	}

	// 创建default Queue 方便用于测试
	defaultQueue := schedulingqueue.NewPrioritySchedulingQueue(fw, "default", "fifo")
	mq.queueMap["default"] = defaultQueue

	return mq, nil
}

func (mq *multiSchedulingQueue) Add(q *v1alpha1.Queue) error {
	pq := schedulingqueue.NewPrioritySchedulingQueue(mq.fw, q.Name, "priority")
	mq.queueMap[pq.Name()] = pq
	return nil
}

func (mq *multiSchedulingQueue) Delete(q *v1alpha1.Queue) error {
	delete(mq.queueMap, q.Name)
	return nil
}

func (mq *multiSchedulingQueue) Update(old *v1alpha1.Queue, new *v1alpha1.Queue) error {
	pq := schedulingqueue.NewPrioritySchedulingQueue(mq.fw, new.Name, "priority")
	mq.queueMap[pq.Name()] = pq
	return nil
}

func (mq *multiSchedulingQueue) GetQueueByName(name string) (queue.SchedulingQueue, bool) {
	if name == "" {
		return mq.queueMap["default"], true
	}
	q, ok := mq.queueMap[name]
	return q, ok
}

func (mq *multiSchedulingQueue) SortedQueue() []queue.SchedulingQueue {
	len := len(mq.queueMap)
	unSortedQueue := make([]queue.SchedulingQueue, len)

	index := 0
	for _, q := range mq.queueMap {
		unSortedQueue[index] = q
		index++
	}

	sort.Slice(unSortedQueue, func(i, j int) bool {
		return mq.lessFunc(unSortedQueue[i].QueueInfo(), unSortedQueue[j].QueueInfo())
	})

	return unSortedQueue
}

func queueInfoKeyFunc(obj interface{}) (string, error) {
	q := obj.(queue.SchedulingQueue)
	return q.Name(), nil
}
