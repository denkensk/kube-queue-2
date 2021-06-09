package extensionclient

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/kube-queue/kube-queue/pkg/apis/queue/v1alpha1"
	"golang.org/x/net/context"
)

type ExtensionClient struct {
	core       http.Client
	ctx        context.Context
	addrPrefix string
}

func (e *ExtensionClient) Close() error {
	return e.Close()
}

func MakeExtensionClient(addr string) (ExtensionClient, error) {
	fileAddr := strings.TrimPrefix(addr, "unix://")
	c := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", fileAddr)
			},
		},
	}

	return ExtensionClient{
		core:       c,
		ctx:        context.Background(),
		addrPrefix: "http://localhost",
	}, nil
}

func (e *ExtensionClient) DequeueJob(qu *v1alpha1.QueueUnit) error {
	data, err := json.Marshal(qu)
	if err != nil {
		return err
	}

	releaseRequest, err := http.NewRequest(
		http.MethodDelete, fmt.Sprintf("%s/", e.addrPrefix), strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	releaseRequest.Header.Set("Content-type", "application/json")

	_, err = e.core.Do(releaseRequest)
	return err
}
