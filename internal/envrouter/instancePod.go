package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	v1 "k8s.io/api/core/v1"
)

type InstancePodService interface {
	FindAll() ([]*api.InstancePod, error)
}

type instancePodService struct {
	podService k8s.PodService
}

func NewInstancePodService(
	podService k8s.PodService,
) InstancePodService {
	return &instancePodService{
		podService,
	}
}

func (i *instancePodService) FindAll() ([]*api.InstancePod, error) {
	var err error
	deployments, err := i.podService.GetAllByLabelExists(k8s.ApplicationLabelKey)
	if err != nil {
		return nil, err
	}
	var result []*api.InstancePod
	for _, v := range deployments {
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
