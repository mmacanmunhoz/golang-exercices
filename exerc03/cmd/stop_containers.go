package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var stopContainers = &cobra.Command{
	Use:   "stop",
	Short: "Para todos os containers em execução",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		defer apiClient.Close()

		containers, err := apiClient.ContainerList(ctx, container.ListOptions{})
		if err != nil {
			fmt.Printf("Erro ao listar containers: %v\n", err)
			return
		}

		if len(containers) == 0 {
			fmt.Println("Nenhum container em execução encontrado.")
			return
		}

		fmt.Printf("Parando %d container(s)...\n", len(containers))
		for _, ctr := range containers {
			fmt.Printf("Parando container %s (%s)... ", ctr.ID[:12], ctr.Names[0])
			timeout := 10
			if err := apiClient.ContainerStop(ctx, ctr.ID, container.StopOptions{Timeout: &timeout}); err != nil {
				fmt.Printf("Erro: %v\n", err)
				continue
			}
			fmt.Println("✓ Parado com sucesso")
		}

	}}

func init() {
	containerCmd.AddCommand(stopContainers)
}
