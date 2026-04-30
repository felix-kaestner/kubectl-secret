package cmd

import (
	"fmt"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type BaseOptions struct {
	configFlags *genericclioptions.ConfigFlags
	streams     genericiooptions.IOStreams

	Name      string
	Namespace string
	Client    typedcorev1.SecretInterface
}

func (o *BaseOptions) Complete(args []string) error {
	if len(args) == 1 {
		o.Name = args[0]
	}

	var err error
	o.Namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return fmt.Errorf("resolving namespace: %w", err)
	}

	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("building REST config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("building client: %w", err)
	}

	o.Client = clientset.CoreV1().Secrets(o.Namespace)
	return nil
}

func (o *BaseOptions) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("secret name is required")
	}
	if o.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	return nil
}
