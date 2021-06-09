package runtime

import (
	"github.com/kube-queue/kube-queue/pkg/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

// PluginFactory is a function that builds a plugin.
type PluginFactory = func(configuration runtime.Object, handle framework.Handle) (framework.Plugin, error)

type Registry map[string]PluginFactory
