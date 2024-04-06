package main

import (
	"fmt"
	"os"

	"github.com/romainframe/grunter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "grunter",
	Version:       grunter.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Grunter is a tool to generate Terragrunt configurations",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
