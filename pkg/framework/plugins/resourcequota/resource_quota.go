package resourcequota

import (
	"context"

	"github.com/kube-queue/kube-queue/pkg/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const Name = "ResourceQuota"

// ResourceQuota is a plugin that implements ResourceQuota filter.
type ResourceQuota struct{}

var _ framework.FilterPlugin = &ResourceQuota{}

// Name returns name of the plugin.
func (rq *ResourceQuota) Name() string {
	return Name
}

func (rq *ResourceQuota) Filter(ctx context.Context, QueueUnit *framework.QueueUnitInfo) *framework.Status {
	return framework.NewStatus(0, "")
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	return &ResourceQuota{}, nil
}
