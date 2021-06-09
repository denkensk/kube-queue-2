package controller

import (
	"github.com/kube-queue/kube-queue/pkg/common/queueserver"
	"github.com/kube-queue/kube-queue/pkg/framework/plugins"
	"github.com/kube-queue/kube-queue/pkg/queue/multischedulingqueue"

	"github.com/kube-queue/kube-queue/pkg/framework/runtime"

	"github.com/kube-queue/kube-queue/pkg/scheduler"

	"github.com/kube-queue/kube-queue/pkg/framework"

	"github.com/kube-queue/kube-queue/pkg/queue"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"

	v1alpha1 "github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
	"github.com/kube-queue/kube-queue/pkg/permission"
	"k8s.io/klog/v2"
)

type Controller struct {
	// extCh is an ExtensionClients worker channel for calling RPC asynchronously
	extCh chan *v1alpha1.QueueUnit

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	multiSchedulingQueue queue.MultiSchedulingQueue
	queueServer          *queueserver.QueueServer
	fw                   *framework.Framework
	scheduler            *scheduler.Scheduler
}

func NewController(
	kubeclientset kubernetes.Interface,
	pc permission.CounterInterface,
	listen string,
	kubeConfigPath string,
	typeAddr map[string]string) (*Controller, error) {
	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})

	schemeModified := scheme.Scheme
	recorder := eventBroadcaster.NewRecorder(schemeModified, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		recorder: recorder,
		extCh:    make(chan *v1alpha1.QueueUnit),
	}
	r := plugins.NewInTreeRegistry()
	fw, err := runtime.NewFramework(r, kubeConfigPath)
	if err != nil {
		klog.Errorf("%s", err)
	}

	multiSchedulingQueue, err := multischedulingqueue.NewMultiSchedulingQueue(fw)
	if err != nil {
		klog.Fatalf("init multi scheduling queue failed %s", err)
	}

	klog.Info("****multiSchedulingQueue")
	controller.queueServer, err = queueserver.NewQueueServer(multiSchedulingQueue, listen, typeAddr)
	if err != nil {
		klog.Fatalf("init queue server failed %s", err)
	}

	klog.Info("****NewQueueServer")
	controller.scheduler, err = scheduler.NewScheduler(multiSchedulingQueue, fw)
	if err != nil {
		klog.Fatalf("init scheduler failed %s", err)
	}

	klog.Info("****NewScheduler")
	return controller, nil
}

func (c *Controller) Start() {
	go c.queueServer.Start()
	c.scheduler.Start()
}
