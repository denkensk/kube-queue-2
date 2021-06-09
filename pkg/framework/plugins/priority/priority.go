package priority

import (
	"github.com/kube-queue/kube-queue/pkg/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const Name = "Priority"

// ResourceQuota is a plugin that implements ResourceQuota filter.
type Priority struct{}

var _ framework.MultiQueueSortPlugin = &Priority{}
var _ framework.QueueSortPlugin = &Priority{}

// Name returns name of the plugin.
func (p *Priority) Name() string {
	return Name
}

func (p *Priority) MultiQueueLess(q1 *framework.QueueInfo, q2 *framework.QueueInfo) bool {
	p1 := q1.Priority()
	p2 := q2.Priority()
	return p1 > p2
}

func (p *Priority) QueueLess(u1 *framework.QueueUnitInfo, u2 *framework.QueueUnitInfo) bool {
	p1 := u1.Unit().Spec.Priority
	p2 := u1.Unit().Spec.Priority
	return p1 > p2
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &Priority{}, nil
}
