package main

import (
	"os"

	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"github.com/felix-kaestner/kubectl-secret/pkg/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-secret", pflag.ExitOnError)
	pflag.CommandLine = flags

	streams := genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	root := cmd.NewCmdSecret(streams)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
