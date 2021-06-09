package runtime

import (
	"github.com/kube-queue/kube-queue/pkg/framework"
	"k8s.io/klog/v2"
)

var _ framework.Framework = &frameworkImpl{}

type frameworkImpl struct {
	multiQueueSortPlugin framework.MultiQueueSortPlugin
	filterPlugins        []framework.FilterPlugin
	queueSortPlugin      []framework.QueueSortPlugin
	kubeConfigPath       string
}

func (f *frameworkImpl) MultiQueueSortFunc() framework.MultiQueueLessFunc {
	return f.multiQueueSortPlugin.MultiQueueLess
}

func (f *frameworkImpl) QueueSortFuncMap() map[string]framework.QueueLessFunc {
	queueLessFuncMap := make(map[string]framework.QueueLessFunc)
	for _, plugin := range f.queueSortPlugin {
		queueLessFuncMap[plugin.Name()] = plugin.QueueLess
	}

	return queueLessFuncMap
}

func (f *frameworkImpl) RunFilterPlugins(unit *framework.QueueUnitInfo) *framework.Status {
	for _, pl := range f.filterPlugins {
		pluginStatus := pl.Filter(nil, unit)
		if pluginStatus.Code() != framework.Success {
			return pluginStatus
		}
	}

	return framework.NewStatus(framework.Success, "")
}

func (f *frameworkImpl) RunScorePlugins() (int64, bool) {
	return 0, false
}

func NewFramework(r Registry, kubeConfigPath string) (framework.Framework, error) {
	filterPlugins := make([]framework.FilterPlugin, 0)
	queueSortPlugin := make([]framework.QueueSortPlugin, 0)
	var multiQueueSortPlugin framework.MultiQueueSortPlugin

	for name, f := range r {
		klog.Infof("init plugins %v", name)
		p, err := f(nil, nil)
		if err != nil {
			klog.Fatalf("init plugin failed %v %v", name, err)
		}
		if i, ok := p.(framework.QueueSortPlugin); ok {
			queueSortPlugin = append(queueSortPlugin, i)
		}
		if i, ok := p.(framework.MultiQueueSortPlugin); ok {
			multiQueueSortPlugin = i
		}
		if i, ok := p.(framework.FilterPlugin); ok {
			filterPlugins = append(filterPlugins, i)
		}
	}
	f := &frameworkImpl{
		kubeConfigPath:       kubeConfigPath,
		filterPlugins:        filterPlugins,
		queueSortPlugin:      queueSortPlugin,
		multiQueueSortPlugin: multiQueueSortPlugin,
	}
	return f, nil
}
