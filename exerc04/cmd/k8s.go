package cmd

import (
	"github.com/spf13/cobra"
)

var containerCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Comandos relacionados a containers Docker",
}

func init() {
	rootCmd.AddCommand(containerCmd)
}
