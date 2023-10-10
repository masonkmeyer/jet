package cmd

import (
	"fmt"
	"os"

	"github.com/masonkmeyer/jet/jet/runtime"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jet",
	Long:  "This tool helps you manage git branches.",
	Short: "Jet is a CLI tool for managing your git branches",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runtime.Run()
	},
}

// Execute is the entry point for the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
