// SPDX-FileCopyrightText: 2026 Felix Kästner
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
)

func NewCmdSecret(streams genericiooptions.IOStreams, version string) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:          "kubectl-secret",
		Short:        "Work with Kubernetes secrets more easily",
		SilenceUsage: true,
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(NewCmdView(configFlags, streams))
	cmd.AddCommand(NewCmdEdit(configFlags, streams))
	cmd.AddCommand(NewCmdVersion(streams, version))

	return cmd
}
