package k8s

import "context"

type DataStorageFactory interface {
	NewRepositoryStorage() ConfigMapDataStorage
	NewApplicationStorage() ConfigMapDataStorage
	NewCredentialsSecretStorage() SecretDataStorage
	NewRefBindingStorage() ConfigMapDataStorage
}

type dataStorageFactory struct {
	client *client
	ns     string
}

func NewDataStorageFactory(client *client) DataStorageFactory {
	ns, err := client.GetK8sNamespace()
	if err != nil {
		panic(err)
	}
	return &dataStorageFactory{client, ns}
}

func (d *dataStorageFactory) NewRepositoryStorage() ConfigMapDataStorage {
	return NewConfigMapDataStorage(
		context.TODO(),
		d.client,
		d.ns,
		RepositoryConfigMapName,
	)
}

func (d *dataStorageFactory) NewApplicationStorage() ConfigMapDataStorage {
	return NewConfigMapDataStorage(
		context.TODO(),
		d.client,
		d.ns,
		ApplicationConfigMapName,
	)
}

func (d *dataStorageFactory) NewCredentialsSecretStorage() SecretDataStorage {
	return NewSecretDataStorage(
		context.TODO(),
		d.client,
		d.ns,
		CredentialsSecretTypeLabelValue,
	)
}

func (d *dataStorageFactory) NewRefBindingStorage() ConfigMapDataStorage {
	return NewConfigMapDataStorage(
		context.TODO(),
		d.client,
		d.ns,
		RefBindingConfigMapName,
	)
}
