package queueserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
	"github.com/kube-queue/kube-queue/pkg/common/extensionclient"
	"github.com/kube-queue/kube-queue/pkg/common/utils"
	"github.com/kube-queue/kube-queue/pkg/queue"
	"golang.org/x/net/context"
	"k8s.io/klog/v2"
)

type QueueServer struct {
	multiSchedulingQueue queue.MultiSchedulingQueue
	extensionClients     map[string]extensionclient.ExtensionClient
	listen               net.Listener
	server               *http.Server
}

func NewQueueServer(multiSchedulingQueue queue.MultiSchedulingQueue, endpoint string, addrs map[string]string) (*QueueServer, error) {
	protocol, addr, err := utils.ParseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	listen, err := net.Listen(protocol, addr)
	if err != nil {
		klog.Info("%v", err)
		return nil, err
	}
	//defer listen.Close()

	s := &QueueServer{
		multiSchedulingQueue: multiSchedulingQueue,
		listen:               listen,
	}

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Error(
					w,
					fmt.Sprintf("Path %s not found", r.URL.Path),
					http.StatusNotFound)
				return
			}

			data, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var qu v1alpha1.QueueUnit
			err = json.Unmarshal(data, &qu)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			switch r.Method {
			case http.MethodPost:
				err = s.add(context.Background(), &qu)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case http.MethodDelete:
				err = s.delete(context.Background(), &qu)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			default:
				http.Error(
					w,
					fmt.Sprintf("method %s not found", r.URL.Path),
					http.StatusNotFound)
				return
			}
		})

	s.server = &http.Server{
		Handler: nil,
	}

	extensionClients := map[string]extensionclient.ExtensionClient{}
	for uType, addr := range addrs {
		ec, err := extensionclient.MakeExtensionClient(addr)
		if err != nil {
			return nil, err
		}
		extensionClients[uType] = ec
	}
	s.extensionClients = extensionClients

	return s, nil
}

func (s *QueueServer) add(ctx context.Context, u *v1alpha1.QueueUnit) error {
	queueName := u.Spec.Queue
	q, ok := s.multiSchedulingQueue.GetQueueByName(queueName)
	if !ok {
		klog.Errorf("queue is not exist %s", queueName)
	}

	err := q.Add(u)
	if err != nil {
		klog.Errorf("queue %s add unit fail", queueName)
	}

	return nil
}

func (s *QueueServer) update(ctx context.Context, u *v1alpha1.QueueUnit) error {
	queueName := u.Spec.Queue
	q, ok := s.multiSchedulingQueue.GetQueueByName(queueName)
	if !ok {
		klog.Errorf("queue is not exist %s", queueName)
	}

	err := q.Delete(u)
	if err != nil {
		klog.Errorf("queue %s add unit fail", queueName)
	}

	return nil
}

func (s *QueueServer) delete(ctx context.Context, u *v1alpha1.QueueUnit) error {
	queueName := u.Spec.Queue
	q, ok := s.multiSchedulingQueue.GetQueueByName(queueName)
	if !ok {
		klog.Errorf("queue is not exist %s", queueName)
	}

	err := q.Delete(u)
	if err != nil {
		klog.Errorf("queue %s add unit fail", queueName)
	}

	return nil
}

func (s *QueueServer) Start() {
	go func() {
		defer s.server.Close()
		if err := s.server.Serve(s.listen); err != nil {
			klog.Fatal("QueueServer error %s", err)
		}
	}()
}

func (s *QueueServer) Dequeue(u *v1alpha1.QueueUnit) error {
	uType := u.Spec.JobType
	ec, ok := s.extensionClients[uType]
	if !ok {
		return fmt.Errorf("extension clients can not be found for %s", uType)
	}

	err := ec.DequeueJob(u)
	if err != nil {
		return err
	}
	return nil
}
