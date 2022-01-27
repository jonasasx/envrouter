package envrouter

import (
	"github.com/ghodss/yaml"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
)

type RepositoryService interface {
	Save(repository *api.Repository) error
	FindByName(repositoryName string) (*api.Repository, error)
	FindAll() ([]*api.Repository, error)
	DeleteByName(repositoryName string) error
	ExistsByName(repositoryName string) bool
}

type repositoryService struct {
	dataStorage k8s.ConfigMapDataStorage
}

func NewRepositoryService(
	dataStorage k8s.ConfigMapDataStorage,
) RepositoryService {
	return &repositoryService{dataStorage: dataStorage}
}

func (r *repositoryService) Save(repository *api.Repository) error {
	value, err := yaml.Marshal(repository)
	if err != nil {
		return err
	}
	return r.dataStorage.Save(repository.Name, string(value))
}

func (r *repositoryService) FindByName(repositoryName string) (*api.Repository, error) {
	data, err := r.dataStorage.GetByKey(repositoryName)
	if err != nil {
		return nil, err
	}
	item := api.Repository{}
	err = yaml.Unmarshal([]byte(data), &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *repositoryService) FindAll() ([]*api.Repository, error) {
	data, err := r.dataStorage.GetAll()
	if err != nil {
		return nil, err
	}
	var result []*api.Repository
	for _, v := range data {
		item := api.Repository{}
		err := yaml.Unmarshal([]byte(v), &item)
		if err != nil {
			return nil, err
		}
		result = append(result, &item)
	}
	return result, nil
}

func (r *repositoryService) DeleteByName(repositoryName string) error {
	return r.dataStorage.DeleteByKey(repositoryName)
}

func (r *repositoryService) ExistsByName(repositoryName string) bool {
	result, err := r.FindByName(repositoryName)
	if err != nil {
		return false
	}
	return result != nil
}
