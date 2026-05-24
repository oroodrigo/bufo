package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bufo",
	Short: "Bufo is a proxy for localhost projects.",
	Long:  "Bufo is a study project to learn how to create a reverse proxy for localhost projects. So, you can run your projects locally and access them through a unique URL.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
