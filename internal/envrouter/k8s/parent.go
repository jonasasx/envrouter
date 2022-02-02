package k8s

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ParentService interface {
	GetPodParents(pod *corev1.Pod) ([]string, error)
}

type parentService struct {
	ctx               context.Context
	client            *client
	replicaSetService ReplicaSetService
}

func NewParentService(
	ctx context.Context,
	client *client,
	replicaSetService ReplicaSetService,
) ParentService {
	return &parentService{
		ctx,
		client,
		replicaSetService,
	}
}

func (p *parentService) GetPodParents(pod *corev1.Pod) ([]string, error) {
	return p.getParents(pod.ObjectMeta)
}

func (p *parentService) getParents(meta metav1.ObjectMeta) ([]string, error) {
	result := []string{}
	for _, v := range meta.OwnerReferences {
		reference := v.APIVersion + "/" + v.Kind
		name := v.Name
		if reference == "apps/v1/ReplicaSet" {
			replicaSet, err := p.replicaSetService.Get(meta.Namespace, name)
			if err != nil {
				return nil, err
			}
			if replicaSet != nil {
				parents, err := p.getParents(replicaSet.ObjectMeta)
				if err != nil {
					return nil, err
				}
				result = append(result, parents...)
			}
		}
		result = append(result, reference+"/"+name)
	}
	return result, nil
}
