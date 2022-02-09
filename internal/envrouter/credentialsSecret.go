package envrouter

import (
	"fmt"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/api"
	"gitlab.com/jonasasx/envrouter/internal/envrouter/k8s"
	rand "gitlab.com/jonasasx/envrouter/internal/utils"
)

type CredentialsSecretService interface {
	Save(credentialsSecretRequest *api.CredentialsSecretRequest) (*api.CredentialsSecretListItem, error)
	FindAll() ([]*api.CredentialsSecretListItem, error)
	FindByName(credentialsSecretName string) (*api.CredentialsSecretRequest, error)
	DeleteByName(name string) error
}

type credentialsSecretService struct {
	dataStorage k8s.SecretDataStorage
}

func NewCredentialsSecretService(
	dataStorage k8s.SecretDataStorage,
) CredentialsSecretService {
	return &credentialsSecretService{dataStorage: dataStorage}
}

func (c *credentialsSecretService) FindAll() ([]*api.CredentialsSecretListItem, error) {
	items, err := c.dataStorage.ListByLabel()
	if err != nil {
		return nil, err
	}
	result := []*api.CredentialsSecretListItem{}
	for k, _ := range items {
		item := api.CredentialsSecretListItem{
			Name: k,
			Type: k8s.CredentialsSecretTypeLabelValue,
		}
		result = append(result, &item)
	}
	return result, nil
}

func (c *credentialsSecretService) FindByName(credentialsSecretName string) (*api.CredentialsSecretRequest, error) {
	item, err := c.dataStorage.GetByName(credentialsSecretName)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("secret %s not found", credentialsSecretName)
	}
	return &api.CredentialsSecretRequest{
		Key:      string(item["key"]),
		Username: string(item["username"]),
		Password: string(item["password"]),
	}, nil
}

func (c *credentialsSecretService) Save(credentialsSecretRequest *api.CredentialsSecretRequest) (*api.CredentialsSecretListItem, error) {
	data := map[string][]byte{
		"key":      []byte(credentialsSecretRequest.Key),
		"username": []byte(credentialsSecretRequest.Username),
		"password": []byte(credentialsSecretRequest.Password),
	}
	name := "envrouter-" + rand.String(8)
	err := c.dataStorage.Save(name, data)
	if err != nil {
		return nil, err
	}
	return &api.CredentialsSecretListItem{
		Name: name,
		Type: k8s.CredentialsSecretTypeLabelValue,
	}, nil
}

func (c *credentialsSecretService) DeleteByName(name string) error {
	return c.dataStorage.DeleteByName(name)
}
