// SPDX-FileCopyrightText: 2026 Felix Kästner
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/felix-kaestner/kubectl-secret/pkg/cmd"
)

// version is set via ldflags at build time.
var version = "dev"

func main() {
	flags := pflag.NewFlagSet("kubectl-secret", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	root := cmd.NewCmdSecret(streams, version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
