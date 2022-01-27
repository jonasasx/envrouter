package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	v1 "k8s.io/api/apps/v1"
)

type InstanceService interface {
	FindAll() ([]*api.Instance, error)
}

type instanceService struct {
	deploymentService k8s.DeploymentService
}

func NewInstanceService(
	deploymentService k8s.DeploymentService,
) InstanceService {
	return &instanceService{
		deploymentService,
	}
}

func (i *instanceService) FindAll() ([]*api.Instance, error) {
	deployments := i.deploymentService.GetAll()
	var result []*api.Instance
	for _, v := range deployments {
		instance := mapInstance(v)
		result = append(result, instance)
	}
	return result, nil
}

func mapInstance(deployment *v1.Deployment) *api.Instance {
	ref := deployment.Annotations[k8s.RefAnnotationKey]
	commitSha := deployment.Annotations[k8s.ShaAnnotationKey]
	return &api.Instance{
		Application: deployment.Labels[k8s.ApplicationLabelKey],
		Ref:         &ref,
		CommitSha:   &commitSha,
		Environment: deployment.Namespace,
		Name:        deployment.Name,
		Type:        "deployment",
	}
}
