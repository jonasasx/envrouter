package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
)

type EnvironmentService interface {
	FindAll() ([]*api.Environment, error)
	ExistsByName(name string) bool
}

type environmentService struct {
	deploymentService k8s.DeploymentService
}

func NewEnvironmentService(
	deploymentService k8s.DeploymentService,
) EnvironmentService {
	return &environmentService{
		deploymentService,
	}
}

func (e *environmentService) FindAll() ([]*api.Environment, error) {
	var err error
	deployments, err := e.deploymentService.GetAllByLabelExists(k8s.ApplicationLabelKey)
	if err != nil {
		return nil, err
	}
	namespaces := map[string]bool{}
	for _, v := range deployments {
		namespaces[v.Namespace] = true
	}
	var result []*api.Environment
	for k, _ := range namespaces {
		environment := api.Environment{Name: k}
		result = append(result, &environment)
	}
	return result, nil
}

func (e *environmentService) ExistsByName(name string) bool {
	deployments, err := e.deploymentService.GetAllInNamespaceByLabelExists(name, k8s.ApplicationLabelKey)
	if err != nil {
		return false
	}
	return len(deployments) > 0

}
