package k8s

import (
	"context"
	"gitlab.com/jonasasx/envrouter/internal/utils"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"time"
)

type DeploymentService interface {
	GetAll() []*v1.Deployment
	GetAllInNamespace(ns string) []*v1.Deployment
	GetAllByLabel(labelName string, labelValue string) []*v1.Deployment
}

type deploymentService struct {
	ctx    context.Context
	client *client
	store  cache.Store
}

func NewDeploymentService(
	ctx context.Context,
	client *client,
	observer utils.Observer,
) (DeploymentService, chan struct{}) {
	var err error
	clientset, _, err := client.getK8sClient()
	if err != nil {
		panic(err)
	}
	optionsModifier := func(options *metav1.ListOptions) {
		options.LabelSelector = ApplicationLabelKey
	}
	watchlist := cache.NewFilteredListWatchFromClient(clientset.AppsV1().RESTClient(), "deployments", "", optionsModifier)
	store, controller := cache.NewInformer(
		watchlist,
		&v1.Deployment{},
		time.Minute*5,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				observer.Publish(nil, obj.(*v1.Deployment))
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				observer.Publish(oldObj.(*v1.Deployment), newObj.(*v1.Deployment))
			},
			DeleteFunc: func(obj interface{}) {
				observer.Publish(obj.(*v1.Deployment), nil)
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	return &deploymentService{
		ctx,
		client,
		store,
	}, stop
}

func (d *deploymentService) GetAllInNamespace(ns string) []*v1.Deployment {
	var result []*v1.Deployment
	deployments := d.store.List()
	for _, v := range deployments {
		deployment := v.(*v1.Deployment)
		if ns == "" || ns == deployment.Namespace {
			result = append(result, deployment)
		}
	}
	return result
}

func (d *deploymentService) GetAll() []*v1.Deployment {
	return d.GetAllInNamespace("")
}

func (d *deploymentService) GetAllByLabel(labelName string, labelValue string) []*v1.Deployment {
	var result []*v1.Deployment
	deployments := d.store.List()
	for _, v := range deployments {
		deployment := v.(*v1.Deployment)
		if _, ok := deployment.Labels[labelName]; ok && deployment.Labels[labelName] == labelValue {
			result = append(result, deployment)
		}
	}
	return result
}
