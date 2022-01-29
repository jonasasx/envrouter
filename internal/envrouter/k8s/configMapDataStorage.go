package k8s

import (
	"context"
	"errors"
	apiv1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigMapDataStorage interface {
	Save(key string, data string) error
	GetAll() (map[string]string, error)
	GetByKey(key string) (string, error)
	DeleteByKey(key string) error
}

type configMapDataStorage struct {
	ctx          context.Context
	client       *client
	namespace    string
	resourceName string
}

func NewConfigMapDataStorage(
	ctx context.Context,
	client *client,
	namespace string,
	resourceName string,
) ConfigMapDataStorage {
	return &configMapDataStorage{
		ctx:          ctx,
		client:       client,
		namespace:    namespace,
		resourceName: resourceName,
	}
}

func (d *configMapDataStorage) Save(key string, data string) error {
	var err error
	clientset, _, err := d.client.getK8sClient()
	if err != nil {
		return nil
	}
	var new bool
	configmap, err := clientset.CoreV1().ConfigMaps(d.namespace).Get(d.ctx, d.resourceName, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
		new = true
	}
	if new {
		configmap = &apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      d.resourceName,
				Namespace: d.namespace,
			},
		}
		configmap.Data = make(map[string]string)
	}
	if configmap.Data == nil {
		configmap.Data = make(map[string]string)
	}

	configmap.Data[key] = data

	if new {
		_, err = clientset.CoreV1().ConfigMaps(d.namespace).Create(d.ctx, configmap, metav1.CreateOptions{})
	} else {
		_, err = clientset.CoreV1().ConfigMaps(d.namespace).Update(d.ctx, configmap, metav1.UpdateOptions{})
	}

	return err
}

func (d *configMapDataStorage) GetAll() (map[string]string, error) {
	var err error
	clientset, _, err := d.client.getK8sClient()
	configmap, err := clientset.CoreV1().ConfigMaps(d.namespace).Get(d.ctx, d.resourceName, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return nil, err
		}
		return make(map[string]string, 0), nil
	}
	return configmap.Data, nil
}

func (d *configMapDataStorage) GetByKey(key string) (string, error) {
	all, err := d.GetAll()
	if err != nil {
		return "", err
	}
	if val, ok := all[key]; ok {
		return val, err
	}
	return "", errors.New("Key " + key + " is not found")
}

func (d *configMapDataStorage) DeleteByKey(key string) error {
	var err error
	clientset, _, err := d.client.getK8sClient()
	configmap, err := clientset.CoreV1().ConfigMaps(d.namespace).Get(d.ctx, d.resourceName, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
		return nil
	}
	if _, ok := configmap.Data[key]; ok {
		delete(configmap.Data, key)
		_, err = clientset.CoreV1().ConfigMaps(d.namespace).Update(d.ctx, configmap, metav1.UpdateOptions{})
		return err
	}
	return nil
}
