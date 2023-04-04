package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	"gitlab.com/jonasasx/envrouter/internal/utils"
	v1 "k8s.io/api/core/v1"
	"reflect"
)

type InstancePodService interface {
	FindAll() ([]*api.InstancePod, error)
}

type instancePodService struct {
	podService    k8s.PodService
	parentService k8s.ParentService
}

func NewInstancePodService(
	podService k8s.PodService,
	instancePodObserver utils.Observer,
	parentService k8s.ParentService,
	podObserver utils.Observer,
) (InstancePodService, chan struct{}) {
	service := &instancePodService{
		podService,
		parentService,
	}
	handler := utils.ObserverEventHandlerFuncs{
		EventFunc: func(oldObj interface{}, newObj interface{}) {
			if oldObj == nil && newObj != nil {
				pod, err := service.mapInstancePod(newObj.(*v1.Pod))
				if err != nil {
					return
				}
				instancePodObserver.Publish(nil, api.SSEvent{
					ItemType: "InstancePod",
					Item:     pod,
					Event:    "UPDATED",
				})
			} else if oldObj != nil && newObj != nil {
				oldPod, err := service.mapInstancePod(oldObj.(*v1.Pod))
				if err != nil {
					return
				}
				newPod, err := service.mapInstancePod(newObj.(*v1.Pod))
				if err != nil {
					return
				}
				if !reflect.DeepEqual(oldPod, newPod) {
					instancePodObserver.Publish(nil, api.SSEvent{
						ItemType: "InstancePod",
						Item:     newPod,
						Event:    "UPDATED",
					})
				}
			} else if oldObj != nil && newObj == nil {
				pod, err := service.mapInstancePod(oldObj.(*v1.Pod))
				if err != nil {
					return
				}
				instancePodObserver.Publish(nil, api.SSEvent{
					ItemType: "InstancePod",
					Item:     pod,
					Event:    "DELETED",
				})
			}
		},
	}
	podObserver.Subscribe(handler)
	stop := make(chan struct{})
	go func() {
		<-stop
		podObserver.Unsubscribe(handler)
	}()
	return service, stop
}

func (i *instancePodService) FindAll() ([]*api.InstancePod, error) {
	pods := i.podService.GetAll()
	var result []*api.InstancePod
	for _, v := range pods {
		instancePod, err := i.mapInstancePod(v)
		if err != nil {
			return nil, err
		}
		result = append(result, instancePod)
	}
	return result, nil
}

func (i *instancePodService) mapInstancePod(pod *v1.Pod) (*api.InstancePod, error) {
	started := true
	ready := true
	var startTime *string
	if pod.Status.StartTime != nil {
		s := pod.Status.StartTime.Time.String()
		startTime = &s
	}
	for _, v := range pod.Status.ContainerStatuses {
		started = started && v.Started != nil && *v.Started
		ready = ready && v.Ready
	}
	ref := pod.Annotations[k8s.RefAnnotationKey]
	commitSha := pod.Annotations[k8s.ShaAnnotationKey]
	parents, err := i.parentService.GetPodParents(pod)
	if err != nil {
		return nil, err
	}
	var environment string
	if env, ok := pod.Labels[k8s.EnvironmentLabelKey]; ok {
		environment = env
	} else {
		environment = pod.Namespace
	}
	return &api.InstancePod{
		Application: pod.Labels[k8s.ApplicationLabelKey],
		Ref:         &ref,
		CommitSha:   &commitSha,
		CreatedTime: pod.CreationTimestamp.String(),
		Environment: environment,
		Name:        pod.Name,
		Phase:       string(pod.Status.Phase),
		Ready:       ready,
		Started:     started,
		StartedTime: startTime,
		Parents:     &parents,
	}, err
}
