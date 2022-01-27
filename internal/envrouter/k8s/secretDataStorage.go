package k8s

import (
	"context"
	apiv1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretDataStorage interface {
	Save(key string, data map[string][]byte) error
	ListByLabel() (map[string]map[string][]byte, error)
	DeleteByName(name string) error
}

type secretDataStorage struct {
	ctx                  context.Context
	client               *client
	namespace            string
	secretTypeLabelValue string
}

func NewSecretDataStorage(
	ctx context.Context,
	client *client,
	namespace string,
	secretTypeLabelValue string,
) SecretDataStorage {
	return &secretDataStorage{
		ctx:                  ctx,
		client:               client,
		namespace:            namespace,
		secretTypeLabelValue: secretTypeLabelValue,
	}
}
func (s *secretDataStorage) Save(key string, data map[string][]byte) error {
	var err error
	clientset, _, err := s.client.getK8sClient()
	if err != nil {
		return nil
	}
	var new bool
	secret, err := clientset.CoreV1().Secrets(s.namespace).Get(s.ctx, key, metav1.GetOptions{})
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return err
		}
		new = true
	}
	if new {
		secret = &apiv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      key,
				Namespace: s.namespace,
				Labels: map[string]string{
					SecretTypeLabelKey: s.secretTypeLabelValue,
				},
			},
		}
	}

	secret.Data = data

	if new {
		_, err = clientset.CoreV1().Secrets(s.namespace).Create(s.ctx, secret, metav1.CreateOptions{})
	} else {
		_, err = clientset.CoreV1().Secrets(s.namespace).Update(s.ctx, secret, metav1.UpdateOptions{})
	}

	return err
}

func (s *secretDataStorage) ListByLabel() (map[string]map[string][]byte, error) {
	var err error
	clientset, _, err := s.client.getK8sClient()
	list, err := clientset.CoreV1().Secrets(s.namespace).List(s.ctx, metav1.ListOptions{LabelSelector: SecretTypeLabelKey + "=" + s.secretTypeLabelValue})
	if err != nil {
		return nil, err
	}
	result := map[string]map[string][]byte{}
	for _, v := range list.Items {
		result[v.Name] = v.Data
	}
	return result, nil
}

func (s *secretDataStorage) DeleteByName(name string) error {
	var err error
	clientset, _, err := s.client.getK8sClient()
	err = clientset.CoreV1().Secrets(s.namespace).Delete(s.ctx, name, metav1.DeleteOptions{})
	return err
}
