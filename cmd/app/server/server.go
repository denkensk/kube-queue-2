package app

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/kube-queue/kube-queue/cmd/app/options"
	"github.com/kube-queue/kube-queue/pkg/controller"
	"github.com/kube-queue/kube-queue/pkg/permission"
	"gopkg.in/yaml.v2"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/sample-controller/pkg/signals"
)

const (
	apiVersion = "v1alpha1"
)

func Run(opt *options.ServerOption) error {
	klog.Errorf("%+v", apiVersion)
	klog.Infof("%+v", apiVersion)

	stopCh := signals.SetupSignalHandler()

	if len(os.Getenv("KUBECONFIG")) > 0 {
		opt.KubeConfig = os.Getenv("KUBECONFIG")
	}

	cfg, err := clientcmd.BuildConfigFromFlags("", opt.KubeConfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s\n", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s\n", err.Error())
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	quotaList := kubeInformerFactory.Core().V1().ResourceQuotas().Lister()

	// Setup the Permission Counter Client
	pc := permission.MakeResourcePermissionCounter(quotaList)

	// Create Extension Client (for release job)
	data, err := ioutil.ReadFile(opt.ExtensionConfig)
	typeAddr := make(map[string]string)
	err = yaml.Unmarshal(data, &typeAddr)
	if err != nil {
		return err
	}

	qController, err := controller.NewController(kubeClient, pc, opt.ListenTo, opt.KubeConfig, typeAddr)
	if err != nil {
		klog.Fatalln("Error building controller\n")
	}
	kubeInformerFactory.Start(stopCh)
	qController.Start()

	return nil
}
