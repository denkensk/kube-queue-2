package scheduler

import (
	"github.com/kube-queue/kube-queue/pkg/common/queueserver"
	"github.com/kube-queue/kube-queue/pkg/framework"
	"github.com/kube-queue/kube-queue/pkg/queue"
	"k8s.io/klog/v2"
)

type Scheduler struct {
	multiSchedulingQueue queue.MultiSchedulingQueue
	fw                   framework.Framework
	queueServer          queueserver.QueueServer
}

func NewScheduler(multiSchedulingQueue queue.MultiSchedulingQueue, fw framework.Framework) (*Scheduler, error) {
	sche := &Scheduler{
		multiSchedulingQueue: multiSchedulingQueue,
		fw:                   fw,
	}
	return sche, nil
}

func (s *Scheduler) Start() {
	s.internalSchedule()
}

// Internal start scheduling
func (s *Scheduler) internalSchedule() {
	for {
		s.schedule()
	}
}

func (s *Scheduler) schedule() {
	sortedQueue := s.multiSchedulingQueue.SortedQueue()
	for _, q := range sortedQueue {
		if q.Length() > 0 {
			klog.Info("schedule cycle")
			unitInfo, err := q.TopUnit()
			if err != nil {
				klog.Errorf("get topunit err %v", err)
			}
			status := s.fw.RunFilterPlugins(unitInfo)
			if status.Code() == framework.Success {
				klog.Infof("dequeue %v", unitInfo.Name())
				err := s.queueServer.Dequeue(unitInfo.Unit())
				if err != nil {
					// 构建一个临时存储的位置
				}
			}
			// 构建一个临时存储的位置
			return
		}
	}
}
