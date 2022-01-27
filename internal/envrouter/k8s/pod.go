package k8s

import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodService interface {
	GetAllByLabelExists(labelName string) ([]*v1.Pod, error)
}

type podService struct {
	ctx    context.Context
	client *client
}

func NewPodService(
	ctx context.Context,
	client *client,
) PodService {
	return &podService{
		ctx,
		client,
	}
}

func (p *podService) GetAllByLabelExists(labelName string) ([]*v1.Pod, error) {
	var err error
	clientset, _, err := p.client.getK8sClient()
	list, err := clientset.CoreV1().Pods("").List(p.ctx, metav1.ListOptions{LabelSelector: labelName})
	if err != nil {
		return nil, err
	}
	var result []*v1.Pod
	for _, v := range list.Items {
		pod := v
		result = append(result, &pod)
	}
	return result, nil
}
