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
	deployments := e.deploymentService.GetAll()
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
	deployments := e.deploymentService.GetAllInNamespace(name)
	return len(deployments) > 0

}
