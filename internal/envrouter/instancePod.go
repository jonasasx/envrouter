package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/utils"
	v1 "k8s.io/api/core/v1"
	"reflect"
)

type InstancePodService interface {
	FindAll() ([]*api.InstancePod, error)
}

type instancePodService struct {
	podService k8s.PodService
}

func NewInstancePodService(
	podServiceFactoryMethod func(*k8s.PodEventHandler) (k8s.PodService, chan struct{}),
	observer utils.Observer,
) (InstancePodService, chan struct{}) {
	handler := &k8s.PodEventHandler{
		AddFunc: func(obj *v1.Pod) {
			observer.Publish(&utils.ObserverEvent{
				Item:  mapInstancePod(obj),
				Event: "UPDATED",
			})
		},
		UpdateFunc: func(oldObj, newObj *v1.Pod) {
			oldPod := mapInstancePod(oldObj)
			newPod := mapInstancePod(newObj)
			if !reflect.DeepEqual(oldPod, newPod) {
				observer.Publish(&utils.ObserverEvent{
					Item:  newPod,
					Event: "UPDATED",
				})
			}
		},
		DeleteFunc: func(obj *v1.Pod) {
			observer.Publish(&utils.ObserverEvent{
				Item:  mapInstancePod(obj),
				Event: "DELETED",
			})
		},
	}
	podService, stop := podServiceFactoryMethod(handler)
	return &instancePodService{
		podService,
	}, stop
}

func (i *instancePodService) FindAll() ([]*api.InstancePod, error) {
	pods := i.podService.GetAll()
	var result []*api.InstancePod
	for _, v := range pods {
		instancePod := mapInstancePod(v)
		result = append(result, instancePod)
	}
	return result, nil
}

func mapInstancePod(pod *v1.Pod) *api.InstancePod {
	started := true
	ready := true
	var startTime *string
	if pod.Status.StartTime != nil {
		s := pod.Status.StartTime.Time.String()
		startTime = &s
	}
	for _, v := range pod.Status.ContainerStatuses {
		started = started && *v.Started
		ready = ready && v.Ready
	}
	ref := pod.Annotations[k8s.RefAnnotationKey]
	commitSha := pod.Annotations[k8s.ShaAnnotationKey]
	return &api.InstancePod{
		Application: pod.Labels[k8s.ApplicationLabelKey],
		Ref:         &ref,
		CommitSha:   &commitSha,
		CreatedTime: pod.CreationTimestamp.String(),
		Environment: pod.Namespace,
		Name:        pod.Name,
		Phase:       string(pod.Status.Phase),
		Ready:       ready,
		Started:     started,
		StartedTime: startTime,
	}
}
