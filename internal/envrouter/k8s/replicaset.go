package k8s

import (
	"context"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type ReplicaSetService interface {
	Get(namespace string, name string) (*v1.ReplicaSet, error)
}

type replicaSetService struct {
	ctx    context.Context
	client *client
	store  cache.Store
}

func NewReplicaSetService(
	ctx context.Context,
	client *client,
) (ReplicaSetService, chan struct{}) {
	var err error
	clientset, _, err := client.getK8sClient()
	if err != nil {
		panic(err)
	}
	optionsModifier := func(options *metav1.ListOptions) {
		options.LabelSelector = ApplicationLabelKey
	}
	watchlist := cache.NewFilteredListWatchFromClient(clientset.AppsV1().RESTClient(), "replicaSets", "", optionsModifier)
	store, controller := cache.NewInformer(
		watchlist,
		&v1.ReplicaSet{},
		time.Minute*5,
		cache.ResourceEventHandlerFuncs{},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	return &replicaSetService{
		ctx,
		client,
		store,
	}, stop
}

func (d *replicaSetService) Get(namespace string, name string) (*v1.ReplicaSet, error) {
	replicaSets, exists, err := d.store.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}
	return replicaSets.(*v1.ReplicaSet), nil
}
