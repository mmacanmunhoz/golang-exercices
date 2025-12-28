package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/matheusmacan/configparser/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

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
}

func init() {
	rootCmd.AddCommand(parseCmd)
	rootCmd.AddCommand(serverCmd)
	parseCmd.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	serverCmd.Flags().StringVarP(&filePath, "file", "f", "", "Arquivo de configuração (YAML ou JSON)")
	parseCmd.MarkFlagRequired("file")
}
