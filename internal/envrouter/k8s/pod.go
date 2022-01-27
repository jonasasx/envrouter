package k8s

import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type PodService interface {
	GetAll() []*v1.Pod
}

type podService struct {
	ctx    context.Context
	client *client
	store  cache.Store
}

func NewPodService(
	ctx context.Context,
	client *client,
) (PodService, chan struct{}) {
	var err error
	clientset, _, err := client.getK8sClient()
	if err != nil {
		panic(err)
	}
	optionsModifier := func(options *metav1.ListOptions) {
		options.LabelSelector = ApplicationLabelKey
	}
	watchlist := cache.NewFilteredListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", "", optionsModifier)
	store, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	return &podService{
		ctx,
		client,
		store,
	}, stop
}

func (p *podService) GetAll() []*v1.Pod {
	pods := p.store.List()
	var result []*v1.Pod
	for _, pod := range pods {
		result = append(result, pod.(*v1.Pod))
	}
	return result
}
