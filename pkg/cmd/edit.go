package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"sigs.k8s.io/yaml"
)

type EditOptions struct {
	BaseOptions
	EditFn func([]byte) ([]byte, error)
}

func NewEditOptions(configFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *EditOptions {
	return &EditOptions{
		BaseOptions: BaseOptions{
			configFlags: configFlags,
			streams:     streams,
		},
		EditFn: openInEditor,
	}
}

func NewCmdEdit(configFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *cobra.Command {
	o := NewEditOptions(configFlags, streams)

	return &cobra.Command{
		Use:          "edit NAME",
		Short:        "Edit a secret with base64-decoded values in your $EDITOR",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			return o.Run(cmd.Context())
		},
	}
}

func (o *EditOptions) Run(ctx context.Context) error {
	secret, err := o.Client.Get(ctx, o.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("getting secret %q: %w", o.Name, err)
	}

	original, err := marshalDecodedSecret(secret)
	if err != nil {
		return err
	}

	edited, err := o.EditFn([]byte(original))
	if err != nil {
		return fmt.Errorf("editing secret: %w", err)
	}

	if bytes.Equal([]byte(original), edited) {
		_, err = fmt.Fprintln(o.streams.Out, "Edit cancelled, no changes made.")
		return err
	}

	updated, err := applyEditedData(secret, edited)
	if err != nil {
		return err
	}

	_, err = o.Client.Update(ctx, updated, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("updating secret %q: %w", o.Name, err)
	}

	_, err = fmt.Fprintf(o.streams.Out, "secret/%s edited\n", o.Name)
	return err
}

func applyEditedData(original *corev1.Secret, editedYAML []byte) (*corev1.Secret, error) {
	var edited struct {
		Data map[string]string `json:"data"`
	}
	if err := yaml.Unmarshal(editedYAML, &edited); err != nil {
		return nil, fmt.Errorf("parsing edited YAML: %w", err)
	}

	updated := original.DeepCopy()
	updated.Data = make(map[string][]byte, len(edited.Data))
	for k, v := range edited.Data {
		updated.Data[k] = []byte(v)
	}

	return updated, nil
}

func openInEditor(content []byte) ([]byte, error) {
	f, err := os.CreateTemp("", "kubectl-secret-*.yaml")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	if _, err := f.Write(content); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("writing temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("closing temp file: %w", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running editor %q: %w", editor, err)
	}

	return os.ReadFile(f.Name())
}
