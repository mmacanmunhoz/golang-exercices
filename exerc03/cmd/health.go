package cmd

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var healthContainer = &cobra.Command{
	Use:   "health [container_id]",
	Short: "Faz o parse de um arquivo de configuração YAML ou JSON",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		containerName := args[0]
		ctx := context.Background()
		apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		defer apiClient.Close()

		inspect, err := apiClient.ContainerInspect(ctx, containerName)
		if err != nil {
			fmt.Printf("Erro ao inspecionar container %s: %v\n", containerName, err)
			return
		}

		fmt.Printf("=== Status de Saúde do Container %s ===\n", containerName)
		fmt.Printf("Estado: %s\n", inspect.State.Status)
		fmt.Printf("Executando: %t\n", inspect.State.Running)
		fmt.Printf("PID: %d\n", inspect.State.Pid)
		fmt.Printf("Código de Saída: %d\n", inspect.State.ExitCode)
		fmt.Printf("Problema de Memória: %t\n", inspect.State.OOMKilled)

	},
}

func init() {
	containerCmd.AddCommand(healthContainer)
}
