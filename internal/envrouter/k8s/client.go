package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type client struct {
	clientConfig clientcmd.ClientConfig
}

func NewClient(context string) *client {
	return &client{clientConfig: getConfig(context)}
}

func getConfig(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}

func (c *client) GetK8sNamespace() (string, error) {
	namespace, _, err := c.clientConfig.Namespace()
	if err != nil {
		return "", nil
	}

	return namespace, nil
}

func (c *client) getK8sClient() (*kubernetes.Clientset, string, error) {
	namespace, err := c.GetK8sNamespace()

	restConfig, err := c.clientConfig.ClientConfig()
	if err != nil {
		return nil, "", err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, "", err
	}
	return clientset, namespace, nil
}
