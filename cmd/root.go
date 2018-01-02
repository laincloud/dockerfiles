package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	commit1      string
	commit2      string
	registryHost string
	organization string
)

var rootCmd = &cobra.Command{
	Use:   "dockerfiles",
	Short: "A helper binary to build, pull, push and retag docker images based on dockerfiles",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
