package config

import (
	"os/exec"
	"strings"
	"testing"
)

func TestValidParseCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "parse", "--file", "../example_config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	if !strings.Contains(string(output), "Configuração carregada com sucesso") {
		t.Errorf("Saída inesperada: %s", output)
	}
}

func TestValidSeverCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "server", "--file", "../example_config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	if !strings.Contains(string(output), "Configuração carregada com sucesso") {
		t.Errorf("Saída inesperada: %s", output)
	}
}

func TestValidHealthCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "health", "--file", "../example_config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "Health Result:") {
		t.Errorf("Não encontrou 'Health Result:' na saída: %s", outputStr)
	}

	if !strings.Contains(outputStr, `"healthy":`) {
		t.Error("JSON não contém campo 'healthy'")
	}

	if !strings.Contains(outputStr, `"worker_id":`) {
		t.Error("JSON não contém campo 'worker_id'")
	}

	if !strings.Contains(outputStr, `"timestamp":`) {
		t.Error("JSON não contém campo 'timestamp'")
	}

	if !strings.Contains(outputStr, "Worker") && !strings.Contains(outputStr, "finished") {
		t.Error("Não detectou workers terminando")
	}
}

func TestValidResponseTimeCommand(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go", "response", "--file", "../example_config.yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Erro ao executar o comando: %v", err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "Response Result:") {
		t.Errorf("Não encontrou 'Response Result:' na saída: %s", outputStr)
	}

	if !strings.Contains(outputStr, `"isfast":`) {
		t.Error("JSON não contém campo 'isfast'")
	}

	if !strings.Contains(outputStr, `"timestamp":`) {
		t.Error("JSON não contém campo 'timestamp'")
	}
}
