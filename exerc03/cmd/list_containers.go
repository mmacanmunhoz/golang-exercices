package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var listContainers = &cobra.Command{
	Use:   "list",
	Short: "Faz o parse de um arquivo de configuração YAML ou JSON",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		defer apiClient.Close()

		containers, err := apiClient.ContainerList(ctx, container.ListOptions{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Containers Docker encontrados:")
		fmt.Println("=================================")

		for _, item := range containers {
			fmt.Println("Name Container:", item.Names)
			fmt.Println("ID Container:", item.ID)
			fmt.Println("Image Container:", item.Image)
			fmt.Println("Status Container:", item.Status)
			fmt.Println("=================================")
		}
	},
}

func init() {
	containerCmd.AddCommand(listContainers)
}
