package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func NewCmdSecret(streams genericiooptions.IOStreams) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:          "kubectl-secret",
		Short:        "Work with Kubernetes secrets more easily",
		SilenceUsage: true,
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(NewCmdView(configFlags, streams))
	cmd.AddCommand(NewCmdEdit(configFlags, streams))

	return cmd
}
