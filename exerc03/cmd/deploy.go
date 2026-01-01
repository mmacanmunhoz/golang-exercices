package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/config"
	"os"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var filePath string

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Faz o parse de um arquivo de configuração YAML ou JSON",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}

		var wg sync.WaitGroup

		var cfg config.Config
		if yaml.Unmarshal(data, &cfg) == nil || json.Unmarshal(data, &cfg) == nil {
			validateConfig(cfg)
			fmt.Printf("Configuração carregada com sucesso:\n%+v\n", cfg)

			deploys := make(chan config.DeployConfig, len(cfg.Deploy))
			for w := 1; w <= 10; w++ {
				wg.Add(1)
				go AsyncBuildImage(&wg, deploys, w)
			}

			for _, deploy := range cfg.Deploy {
				deploys <- deploy
			}

			close(deploys)
			wg.Wait()

		} else {
			fmt.Println("Erro ao fazer o parse do arquivo de configuração.")
			os.Exit(1)
		}

	},
}

func validateConfig(cfg config.Config) {
	for i, deploy := range cfg.Deploy {
		if deploy.Name == "" || deploy.Image == "" {
			fmt.Printf("Deploy #%d com campos obrigatórios ausentes\n", i)
		}
	}
}

func AsyncBuildImage(wg *sync.WaitGroup, deploys <-chan config.DeployConfig, id int) {
	defer wg.Done()

	ctx := context.Background()
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	defer apiClient.Close()

	shouldReturn := DeployContainers(ctx, deploys, apiClient)
	if shouldReturn {
		return
	}
}

func DeployContainers(ctx context.Context, deploys <-chan config.DeployConfig, apiClient *client.Client) bool {
	for deploy := range deploys {
		imageName := deploy.Image
		out, err := apiClient.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			fmt.Printf("Erro ao fazer pull da imagem %s: %v\n", imageName, err)
			continue
		}
		defer out.Close()
		fmt.Printf("Pull da imagem %s realizado com sucesso\n", imageName)

		io.Copy(os.Stdout, out)

		exposedPorts := make(nat.PortSet)
		portBindings := make(nat.PortMap)

		containerConfig, hostConfig := PortBinding(deploy, exposedPorts, portBindings)

resp, create := CreateContainer(ctx, deploy, apiClient, containerConfig, hostConfig)
		if create {
			return true
		}

		init := InitContainer(ctx, deploy, apiClient, resp)
		if init {
			return true
		}

		NetworkConnection(ctx, deploy, apiClient, resp)

	}
	return false
}

func PortBinding(deploy config.DeployConfig, exposedPorts nat.PortSet, portBindings nat.PortMap) (*container.Config, *container.HostConfig) {
	for containerPort, hostPort := range deploy.Port {
		port, err := nat.NewPort("tcp", strings.Split(containerPort, "/")[0])
		if err != nil {
			fmt.Printf("Erro ao configurar porta %s: %v\n", containerPort, err)
			continue
		}
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{HostPort: hostPort},
		}
	}

	containerConfig := &container.Config{
		Image:        deploy.Image,
		Env:          deploy.Env,
		Cmd:          deploy.Command,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        deploy.Volume,
	}
	return containerConfig, hostConfig
}

func CreateContainer(ctx context.Context, deploy config.DeployConfig, apiClient *client.Client, containerConfig *container.Config, hostConfig *container.HostConfig) (container.CreateResponse, bool) {
	networkingConfig := &network.NetworkingConfig{}
	resp, err := apiClient.ContainerCreate(ctx, containerConfig, hostConfig, networkingConfig, nil, deploy.Name)
	if err != nil {
		fmt.Printf("Erro ao criar container: %v\n", err)
		return container.CreateResponse{}, true
	}
	return resp, false
}

func NetworkConnection(ctx context.Context, deploy config.DeployConfig, apiClient *client.Client, resp container.CreateResponse) {
	for _, networkName := range deploy.Networks {
		if err := apiClient.NetworkConnect(ctx, networkName, resp.ID, nil); err != nil {
			fmt.Printf("Aviso: Erro ao conectar à rede %s: %v\n", networkName, err)
		} else {
			fmt.Printf("✓ Container conectado à rede %s\n", networkName)
		}
	}
}

func InitContainer(ctx context.Context, deploy config.DeployConfig, apiClient *client.Client, resp container.CreateResponse) bool {
	fmt.Printf("Iniciando container %s...\n", deploy.Name)
	if err := apiClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		fmt.Printf("Erro ao iniciar container: %v\n", err)
		return true
	}

	fmt.Printf("✓ Container %s iniciado com sucesso! ID: %s\n", deploy.Name, resp.ID[:12])
	return false
}

func init() {
	containerCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	deployCmd.MarkFlagRequired("file")
}
