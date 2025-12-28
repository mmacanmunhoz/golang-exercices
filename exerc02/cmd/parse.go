package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"configparser-exerc02/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type HealthResult struct {
	config.ServerConfig
	Healthy    bool   `json:"healthy"`
	WorkerID   int    `json:"worker_id"`
	StatusCode int    `json:"status_code"`
	Timestamp  string `json:"timestamp"`
}

type ResponseResult struct {
	config.WebsiteConfig
	Isfast    bool   `json:"isfast"`
	Timestamp string `json:"timestamp"`
}

var filePath string

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Faz o parse de um arquivo de configuração YAML ou JSON",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}

		var cfg config.Config
		if yaml.Unmarshal(data, &cfg) == nil || json.Unmarshal(data, &cfg) == nil {
			validateConfig(cfg)
			fmt.Printf("Configuração carregada com sucesso:\n%+v\n", cfg)
		} else {
			fmt.Println("Erro ao fazer o parse do arquivo de configuração.")
			os.Exit(1)
		}
	},
}

var responseCheck = &cobra.Command{
	Use:   "response",
	Short: "Testa o endpoint de tempo de resposta",
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
			webservers := make(chan config.WebsiteConfig, len(cfg.Website))
			for w := 1; w <= 10; w++ {
				wg.Add(1)
				go AsyncResponseTime(&wg, webservers, w)
			}

			for _, webserver := range cfg.Website {
				webservers <- webserver
			}

			close(webservers)
			wg.Wait()
		}

	},
}

var testHealthStatus = &cobra.Command{
	Use:   "health",
	Short: "Testa o endpoint de health check",
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
			servers := make(chan config.ServerConfig, len(cfg.Servers))
			for w := 1; w <= 10; w++ {
				wg.Add(1)
				go AsyncHealthCheck(&wg, servers, w)
			}

			for _, server := range cfg.Servers {
				servers <- server
			}

			close(servers)
			wg.Wait()
		}
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Imprimi somente os servidores",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}

		var cfg config.Config
		if yaml.Unmarshal(data, &cfg) == nil || json.Unmarshal(data, &cfg) == nil {
			validateConfig(cfg)
			fmt.Printf("Configuração carregada com sucesso:\n%+v\n", cfg.Servers)
		} else {
			fmt.Println("Erro ao fazer o parse do arquivo de configuração.")
			os.Exit(1)
		}
	},
}

func validateConfig(cfg config.Config) {
	for i, server := range cfg.Servers {
		if server.Name == "" || server.Host == "" || server.Port == 0 {
			fmt.Printf("Servidor #%d com campos obrigatórios ausentes\n", i)
		}
	}
	db := cfg.Database
	if db.Host == "" || db.Port == 0 || db.User == "" {
		fmt.Println("Configuração do banco de dados com campos obrigatórios ausentes")
	}

	for i, website := range cfg.Website {
		if website.Name == "" || website.Url == "" || website.MaxResponseTime == 0 {
			fmt.Printf("Website #%d com campos obrigatórios ausentes\n", i)
		}
	}
}

func AsyncResponseTime(wg *sync.WaitGroup, webservers <-chan config.WebsiteConfig, id int) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for webserver := range webservers {
		if webserver.Url != "" {
			start := time.Now()
			req, err := http.NewRequestWithContext(ctx, "GET", webserver.Url, nil)

			if err != nil {
				fmt.Printf("Erro ao acessar o website Worker %d (%s): %v\n", id, webserver.Name, err)
				continue
			}

			duration := time.Since(start)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Erro ao acessar o website Worker %d (%s): %v\n", id, webserver.Url, err)
				continue
			}

			defer resp.Body.Close()

			milliseconds := duration.Microseconds()

			isFAst := milliseconds < int64(webserver.MaxResponseTime)

			result := ResponseResult{
				WebsiteConfig: webserver,
				Isfast:        isFAst,
				Timestamp:     time.Now().Format(time.RFC3339),
			}
			jsonData, _ := json.Marshal(result)
			fmt.Printf("Response Result: %s\n", jsonData)

		}

	}

}

func AsyncHealthCheck(wg *sync.WaitGroup, servers <-chan config.ServerConfig, id int) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for server := range servers {
		if server.Host != "" {

			url := fmt.Sprintf("%s://%s:%d/%s", server.Protocol, server.Host, server.Port, server.Healthcheck)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

			if err != nil {
				fmt.Printf("Erro ao acessar o servidor Worker %d (%s): %v\n", id, server.Name, err)
				continue
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("Erro ao acessar o servidor Worker %d (%s): %v\n", id, server.Name, err)
				continue
			}

			health := resp.StatusCode == http.StatusOK

			defer resp.Body.Close()

			result := HealthResult{
				ServerConfig: server,
				Healthy:      health,
				WorkerID:     id,
				StatusCode:   resp.StatusCode,
				Timestamp:    time.Now().Format(time.RFC3339),
			}
			jsonData, _ := json.Marshal(result)
			fmt.Printf("Health Result: %s\n", jsonData)

		}
	}
	fmt.Printf("Worker %d finished\n", id)
}

func init() {
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(testHealthStatus)
	rootCmd.AddCommand(responseCheck)
	parseCmd.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	serverCmd.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	testHealthStatus.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	responseCheck.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	parseCmd.MarkFlagRequired("file")
	serverCmd.MarkFlagRequired("file")
	testHealthStatus.MarkFlagRequired("file")
	responseCheck.MarkFlagRequired("file")
}
