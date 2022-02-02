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
	podServiceFactoryMethod func(*k8s.PodEventHandler) (k8s.PodService, chan struct{}),
	observer utils.Observer,
	parentService k8s.ParentService,
) (InstancePodService, chan struct{}) {
	handler := &k8s.PodEventHandler{
		AddFunc: func(obj *v1.Pod) {
			pod, err := mapInstancePod(obj, parentService)
			if err != nil {
				return
			}
			observer.Publish(nil, api.SSEvent{
				ItemType: "InstancePod",
				Item:     pod,
				Event:    "UPDATED",
			})
		},
		UpdateFunc: func(oldObj, newObj *v1.Pod) {
			oldPod, err := mapInstancePod(oldObj, parentService)
			if err != nil {
				return
			}
			newPod, err := mapInstancePod(newObj, parentService)
			if err != nil {
				return
			}
			if !reflect.DeepEqual(oldPod, newPod) {
				observer.Publish(nil, api.SSEvent{
					ItemType: "InstancePod",
					Item:     newPod,
					Event:    "UPDATED",
				})
			}
		},
		DeleteFunc: func(obj *v1.Pod) {
			pod, err := mapInstancePod(obj, parentService)
			if err != nil {
				return
			}
			observer.Publish(nil, api.SSEvent{
				ItemType: "InstancePod",
				Item:     pod,
				Event:    "DELETED",
			})
		},
	}
	podService, stop := podServiceFactoryMethod(handler)
	return &instancePodService{
		podService,
		parentService,
	}, stop
}

func (i *instancePodService) FindAll() ([]*api.InstancePod, error) {
	pods := i.podService.GetAll()
	var result []*api.InstancePod
	for _, v := range pods {
		instancePod, err := mapInstancePod(v, i.parentService)
		if err != nil {
			return nil, err
		}
		result = append(result, instancePod)
	}
	return result, nil
}

func mapInstancePod(pod *v1.Pod, parentService k8s.ParentService) (*api.InstancePod, error) {
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
	parents, err := parentService.GetPodParents(pod)
	if err != nil {
		return nil, err
	}
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
		Parents:     &parents,
	}, err
}
