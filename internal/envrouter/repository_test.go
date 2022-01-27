package envrouter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type dataStorageMock struct {
}

func (d *dataStorageMock) Save(key string, data string) error {
	return nil
}

func (d *dataStorageMock) GetAll() (map[string]string, error) {
	return map[string]string{
		"RepoName":   `{"name": "RepoName", "url": "https://envrouter.io", "credentialsSecret": "secret"}`,
		"RepoName 2": `{"name": "RepoName 2", "url": "https://envrouter.io/resource", "credentialsSecret": "secret2"}`,
	}, nil
}

func (d *dataStorageMock) GetByKey(key string) (string, error) {
	return `{"name": "RepoName", "url": "https://envrouter.io", "credentialsSecret": "secret"}`, nil
}

func (d *dataStorageMock) DeleteByKey(key string) error {
	return nil
}

func TestRepository_Save(t *testing.T) {
	service := NewRepositoryService(&dataStorageMock{})
	repository := &Repository{}
	err := service.Save(repository)
	assert.Nil(t, err)
}

func TestRepository_FindByName(t *testing.T) {
	service := NewRepositoryService(&dataStorageMock{})
	repository, err := service.FindByName("RepoName")
	assert.Nil(t, err)
	assert.Equal(t, "RepoName", repository.Name)
	assert.Equal(t, "https://envrouter.io", repository.Url)
	assert.Equal(t, "secret", repository.CredentialsSecret)
}

func TestRepository_FindAll(t *testing.T) {
	service := NewRepositoryService(&dataStorageMock{})
	repositories, err := service.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(repositories))
	assert.Equal(t, "RepoName", repositories[0].Name)
	assert.Equal(t, "https://envrouter.io", repositories[0].Url)
	assert.Equal(t, "secret", repositories[0].CredentialsSecret)
	assert.Equal(t, "RepoName 2", repositories[1].Name)
	assert.Equal(t, "https://envrouter.io/resource", repositories[1].Url)
	assert.Equal(t, "secret2", repositories[1].CredentialsSecret)
}

func TestRepository_Delete(t *testing.T) {
	service := NewRepositoryService(&dataStorageMock{})
	err := service.DeleteByName("RepoName")
	assert.Nil(t, err)
}
