package k8s

import (
	"context"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentService interface {
	GetAllByLabelExists(labelName string) ([]*v1.Deployment, error)
	GetAllInNamespaceByLabelExists(ns string, labelName string) ([]*v1.Deployment, error)
}

type deploymentService struct {
	ctx    context.Context
	client *client
}

func NewDeploymentService(
	ctx context.Context,
	client *client,
) DeploymentService {
	return &deploymentService{
		ctx,
		client,
	}
}

func (d *deploymentService) GetAllInNamespaceByLabelExists(ns string, labelName string) ([]*v1.Deployment, error) {
	var err error
	clientset, _, err := d.client.getK8sClient()
	list, err := clientset.AppsV1().Deployments(ns).List(d.ctx, metav1.ListOptions{LabelSelector: labelName})
	if err != nil {
		return nil, err
	}
	var result []*v1.Deployment
	for _, v := range list.Items {
		deployment := v
		result = append(result, &deployment)
	}
	return result, nil
}

func (d *deploymentService) GetAllByLabelExists(labelName string) ([]*v1.Deployment, error) {
	return d.GetAllInNamespaceByLabelExists("", labelName)
}
