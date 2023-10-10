package cmd

import (
	"fmt"
	"os"

	"github.com/masonkmeyer/jet/jet/runtime"
	"github.com/spf13/cobra"
)

func newRootCmd() (*cobra.Command, chan string) {
	exitChannel := make(chan string, 1)

	return &cobra.Command{
		Use:   "jet",
		Long:  "This tool helps you manage git branches.",
		Short: "Jet is a CLI tool for managing your git branches",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			runtime.Run(exitChannel)
		},
	}, exitChannel
}

// Execute is the entry point for the CLI
func Execute() chan string {
	c, exitChannel := newRootCmd()

	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return exitChannel
}
