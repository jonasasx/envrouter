package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	"gitlab.com/jonasasx/envrouter/internal/utils"
	v1 "k8s.io/api/apps/v1"
	"reflect"
)

type InstanceService interface {
	FindAll() ([]*api.Instance, error)
}

type instanceService struct {
	deploymentService k8s.DeploymentService
}

func NewInstanceService(
	deploymentService k8s.DeploymentService,
	instanceObserver utils.Observer,
	deploymentObserver utils.Observer,
) (InstanceService, chan struct{}) {
	service := &instanceService{
		deploymentService,
	}
	handler := utils.ObserverEventHandlerFuncs{
		EventFunc: func(oldObj interface{}, newObj interface{}) {
			if oldObj == nil && newObj != nil {
				instance := service.mapInstance(newObj.(*v1.Deployment))
				instanceObserver.Publish(nil, api.SSEvent{
					ItemType: "Instance",
					Item:     instance,
					Event:    "UPDATED",
				})
			} else if oldObj != nil && newObj != nil {
				oldInstance := service.mapInstance(oldObj.(*v1.Deployment))
				newInstance := service.mapInstance(newObj.(*v1.Deployment))
				if !reflect.DeepEqual(oldInstance, newInstance) {
					instanceObserver.Publish(nil, api.SSEvent{
						ItemType: "Instance",
						Item:     newInstance,
						Event:    "UPDATED",
					})
				}
			} else if oldObj != nil && newObj == nil {
				instance := service.mapInstance(oldObj.(*v1.Deployment))
				instanceObserver.Publish(nil, api.SSEvent{
					ItemType: "Instance",
					Item:     instance,
					Event:    "DELETED",
				})
			}
		},
	}
	deploymentObserver.Subscribe(handler)
	stop := make(chan struct{})
	go func() {
		<-stop
		deploymentObserver.Unsubscribe(handler)
	}()
	return service, stop
}

func (i *instanceService) FindAll() ([]*api.Instance, error) {
	deployments := i.deploymentService.GetAll()
	var result []*api.Instance
	for _, v := range deployments {
		instance := i.mapInstance(v)
		result = append(result, instance)
	}
	return result, nil
}

func (i *instanceService) mapInstance(deployment *v1.Deployment) *api.Instance {
	ref := deployment.Annotations[k8s.RefAnnotationKey]
	commitSha := deployment.Annotations[k8s.ShaAnnotationKey]
	var environment string
	if env, ok := deployment.Labels[k8s.EnvironmentLabelKey]; ok {
		environment = env
	} else {
		environment = deployment.Namespace
	}
	return &api.Instance{
		Application: deployment.Labels[k8s.ApplicationLabelKey],
		Ref:         &ref,
		CommitSha:   &commitSha,
		Environment: environment,
		Name:        deployment.Name,
		Type:        "apps/v1/Deployment",
	}
}
