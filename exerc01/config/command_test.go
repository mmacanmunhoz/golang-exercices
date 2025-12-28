package config

import (
	"os/exec"
	"strings"
	"testing"
)

func TestValidParseCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "parse", "--file", "../../../example_config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	if !strings.Contains(string(output), "Configuração carregada com sucesso") {
		t.Errorf("Saída inesperada: %s", output)
	}
}

func TestValidSeverCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "server", "--file", "../../../example_config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	if !strings.Contains(string(output), "Configuração carregada com sucesso") {
		t.Errorf("Saída inesperada: %s", output)
	}
}
