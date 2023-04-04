package envrouter

import (
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	"sort"
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
		if env, ok := v.Labels[k8s.EnvironmentLabelKey]; ok {
			namespaces[env] = true
		} else {
			namespaces[v.Namespace] = true
		}
	}
	var result []*api.Environment
	for k, _ := range namespaces {
		environment := api.Environment{Name: k}
		result = append(result, &environment)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result, nil
}

func (e *environmentService) ExistsByName(name string) bool {
	return len(e.deploymentService.GetAllInNamespace(name)) > 0 || len(e.deploymentService.GetAllByLabel(k8s.EnvironmentLabelKey, name)) > 0

}
