package cmd_test

import (
	"bytes"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/felix-kaestner/kubectl-secret/pkg/cmd"
)

func TestEditValidate(t *testing.T) {
	streams := genericiooptions.IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}}

	t.Run("missing name returns error", func(t *testing.T) {
		o := cmd.NewEditOptions(nil, streams)
		o.Namespace = "default"
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for missing name, got nil")
		}
	})

	t.Run("missing namespace returns error", func(t *testing.T) {
		o := cmd.NewEditOptions(nil, streams)
		o.Name = "my-secret"
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for missing namespace, got nil")
		}
	})

	t.Run("name and namespace set is valid", func(t *testing.T) {
		o := cmd.NewEditOptions(nil, streams)
		o.Name = "my-secret"
		o.Namespace = "default"
		if err := o.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestEditNoop(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key": []byte("value"),
		},
	}

	clientset := fake.NewClientset(secret)
	streams := genericiooptions.IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}}

	o := cmd.NewEditOptions(nil, streams)
	o.Name = "my-secret"
	o.Namespace = "default"
	o.Client = clientset.CoreV1().Secrets("default")
	o.EditFn = func(content []byte) ([]byte, error) { return content, nil }

	if err := o.Run(t.Context()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, action := range clientset.Actions() {
		if action.GetVerb() == "patch" || action.GetVerb() == "update" {
			t.Errorf("expected no write actions for no-op edit, got: %s", action.GetVerb())
		}
	}
}

func TestEditAppliesChanges(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key": []byte("original"),
		},
	}

	clientset := fake.NewClientset(secret)
	streams := genericiooptions.IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}}

	o := cmd.NewEditOptions(nil, streams)
	o.Name = "my-secret"
	o.Namespace = "default"
	o.Client = clientset.CoreV1().Secrets("default")
	o.EditFn = func(content []byte) ([]byte, error) {
		return bytes.ReplaceAll(content, []byte("original"), []byte("updated")), nil
	}

	if err := o.Run(t.Context()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	patched := false
	for _, action := range clientset.Actions() {
		if action.GetVerb() == "update" {
			patched = true
		}
	}
	if !patched {
		t.Error("expected an update action after edit, got none")
	}
}
