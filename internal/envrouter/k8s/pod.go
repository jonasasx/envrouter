package k8s

import (
	"context"
	"gitlab.com/jonasasx/envrouter/internal/utils"
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
	observer utils.Observer,
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
		time.Minute*5,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				observer.Publish(nil, obj.(*v1.Pod))
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				observer.Publish(oldObj.(*v1.Pod), newObj.(*v1.Pod))
			},
			DeleteFunc: func(obj interface{}) {
				observer.Publish(obj.(*v1.Pod), nil)
			},
		},
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
