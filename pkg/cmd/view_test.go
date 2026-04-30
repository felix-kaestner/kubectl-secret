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

func TestViewValidate(t *testing.T) {
	streams := genericiooptions.IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}}

	t.Run("missing name returns error", func(t *testing.T) {
		o := cmd.NewViewOptions(nil, streams)
		o.Namespace = "default"
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for missing name, got nil")
		}
	})

	t.Run("missing namespace returns error", func(t *testing.T) {
		o := cmd.NewViewOptions(nil, streams)
		o.Name = "my-secret"
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for missing namespace, got nil")
		}
	})

	t.Run("name and namespace set is valid", func(t *testing.T) {
		o := cmd.NewViewOptions(nil, streams)
		o.Name = "my-secret"
		o.Namespace = "default"
		if err := o.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestViewRun(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("s3cr3t"),
		},
	}

	out := &bytes.Buffer{}
	streams := genericiooptions.IOStreams{Out: out, ErrOut: &bytes.Buffer{}}

	o := cmd.NewViewOptions(nil, streams)
	o.Name = "my-secret"
	o.Namespace = "default"
	o.Client = fake.NewClientset(secret).CoreV1().Secrets("default")

	if err := o.Run(t.Context()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !bytes.Contains([]byte(output), []byte("username")) || !bytes.Contains([]byte(output), []byte("admin")) {
		t.Errorf("expected decoded username in output, got:\n%s", output)
	}
	if !bytes.Contains([]byte(output), []byte("password")) || !bytes.Contains([]byte(output), []byte("s3cr3t")) {
		t.Errorf("expected decoded password in output, got:\n%s", output)
	}
}
