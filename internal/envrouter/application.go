package envrouter

import (
	"errors"
	"github.com/ghodss/yaml"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
)

type ApplicationService interface {
	FindAll() ([]*api.Application, error)
	Save(a *api.Application) (*api.Application, error)
	ExistsByName(name string) bool
}

type applicationService struct {
	deploymentService  k8s.DeploymentService
	applicationStorage k8s.ConfigMapDataStorage
	repositoryService  RepositoryService
}

func NewApplicationService(
	deploymentService k8s.DeploymentService,
	applicationStorage k8s.ConfigMapDataStorage,
	repositoryService RepositoryService,
) ApplicationService {
	return &applicationService{
		deploymentService,
		applicationStorage,
		repositoryService,
	}
}

func (a *applicationService) FindAll() ([]*api.Application, error) {
	var err error
	deployments, err := a.deploymentService.GetAllByLabelExists(k8s.ApplicationLabelKey)
	if err != nil {
		return nil, err
	}
	applicationNames := map[string]bool{}
	for _, v := range deployments {
		if applicationName, ok := v.Labels[k8s.ApplicationLabelKey]; ok {
			applicationNames[applicationName] = true
		}
	}
	applicationConfigs, err := a.applicationStorage.GetAll()
	if err != nil {
		return nil, err
	}
	var result []*api.Application
	for applicationName, _ := range applicationNames {
		var application *api.Application

		if config, ok := applicationConfigs[applicationName]; ok {
			application, err = unmarshallApplicationConfig(config)
			if err != nil {
				return nil, err
			}
			application.Name = applicationName
		} else {
			application = &api.Application{Name: applicationName}
		}

		result = append(result, application)
	}
	return result, nil
}

func (a *applicationService) Save(application *api.Application) (*api.Application, error) {
	if application.RepositoryName != nil && len(*application.RepositoryName) > 0 && !a.repositoryService.ExistsByName(*application.RepositoryName) {
		return nil, errors.New("Repository " + *application.RepositoryName + " is not found")
	}
	applicationConfig, err := marshallApplicationConfig(application)
	if err != nil {
		return nil, err
	}
	err = a.applicationStorage.Save(application.Name, applicationConfig)
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (a *applicationService) ExistsByName(name string) bool {
	var err error
	deployments, err := a.deploymentService.GetAllByLabelExists(k8s.ApplicationLabelKey)
	if err != nil {
		return false
	}
	for _, v := range deployments {
		if applicationName, ok := v.Labels[k8s.ApplicationLabelKey]; ok {
			if applicationName == name {
				return true
			}
		}
	}
	return false
}

func marshallApplicationConfig(application *api.Application) (string, error) {
	result, err := yaml.Marshal(application)
	if err != nil {
		return "", err
	}
	return string(result), err
}

func unmarshallApplicationConfig(s string) (*api.Application, error) {
	item := api.Application{}
	err := yaml.Unmarshal([]byte(s), &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
