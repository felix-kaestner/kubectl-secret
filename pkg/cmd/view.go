package cmd

import (
	"context"
	"fmt"
	"io"
	"maps"
	"slices"
	"text/tabwriter"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

type ViewOptions struct {
	BaseOptions
}

func NewViewOptions(configFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *ViewOptions {
	return &ViewOptions{
		BaseOptions: BaseOptions{
			configFlags: configFlags,
			streams:     streams,
		},
	}
}

func NewCmdView(configFlags *genericclioptions.ConfigFlags, streams genericiooptions.IOStreams) *cobra.Command {
	o := NewViewOptions(configFlags, streams)

	return &cobra.Command{
		Use:          "view NAME",
		Short:        "Display a secret with base64-decoded values",
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

func (o *ViewOptions) Run(ctx context.Context) error {
	secret, err := o.Client.Get(ctx, o.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("getting secret %q: %w", o.Name, err)
	}

	return printDecodedSecret(secret, o.streams.Out)
}

// printDecodedSecret prints a secret in the same style as kubectl describe secret.
// Adapted from https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/kubectl/pkg/describe/describe.go
// The only difference: all data values are shown decoded instead of as byte counts.
func printDecodedSecret(secret *corev1.Secret, out io.Writer) error { //nolint:errcheck
	w := tabwriter.NewWriter(out, 0, 8, 2, ' ', 0)

	fmt.Fprintf(w, "Name:\t%s\n", secret.Name)
	fmt.Fprintf(w, "Namespace:\t%s\n", secret.Namespace)

	if len(secret.Labels) > 0 {
		for i, k := range slices.Sorted(maps.Keys(secret.Labels)) {
			if i == 0 {
				fmt.Fprintf(w, "Labels:\t%s=%s\n", k, secret.Labels[k])
			} else {
				fmt.Fprintf(w, "\t%s=%s\n", k, secret.Labels[k])
			}
		}
	} else {
		fmt.Fprintf(w, "Labels:\t<none>\n")
	}

	delete(secret.Annotations, corev1.LastAppliedConfigAnnotation)
	if len(secret.Annotations) > 0 {
		for i, k := range slices.Sorted(maps.Keys(secret.Annotations)) {
			if i == 0 {
				fmt.Fprintf(w, "Annotations:\t%s=%s\n", k, secret.Annotations[k])
			} else {
				fmt.Fprintf(w, "\t%s=%s\n", k, secret.Annotations[k])
			}
		}
	} else {
		fmt.Fprintf(w, "Annotations:\t<none>\n")
	}

	fmt.Fprintf(w, "\nType:\t%s\n", secret.Type)

	fmt.Fprintf(w, "\nData\n====\n")
	for _, k := range slices.Sorted(maps.Keys(secret.Data)) {
		// Unlike kubectl describe, we always show decoded values instead of byte counts.
		fmt.Fprintf(w, "%s:\t%s\n", k, string(secret.Data[k]))
	}

	return w.Flush()
}
