package cmd

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var listLogs = &cobra.Command{
	Use:   "logs [container_id] ",
	Short: "Faz o parse de um arquivo de configuração YAML ou JSON",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		containerID := args[0]
		ctx := context.Background()
		apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		defer apiClient.Close()

		options := container.LogsOptions{ShowStdout: true, Follow: true, Tail: "10"}
		out, err := apiClient.ContainerLogs(ctx, containerID, options)
		if err != nil {
			panic(err)
		}

		io.Copy(os.Stdout, out)
	},
}

func init() {
	containerCmd.AddCommand(listLogs)
}
