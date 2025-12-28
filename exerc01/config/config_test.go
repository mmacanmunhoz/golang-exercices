package config

import (
	"testing"
)

func TestValidConfig(t *testing.T) {
	cfg := Config{
		Servers: []ServerConfig{{
			Name:     "app-server",
			Host:     "localhost",
			Port:     8080,
			Replicas: 3,
		}},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "admin",
			Password: "secret",
		},
	}

	if cfg.Servers[0].Name == "" {
		t.Error("Nome do servidor deve estar vazio")
	}
	if cfg.Database.Host == "" {
		t.Error("Host do banco de dados deve estar vazio")
	}
}
