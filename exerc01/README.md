# Config Parser CLI

Parser de arquivos de configuração YAML e JSON usando Cobra CLI.

## Descrição

Este projeto é uma ferramenta de linha de comando desenvolvida em Go que permite fazer o parse e validação de arquivos de configuração nos formatos YAML e JSON. O aplicativo utiliza o framework Cobra para criar uma interface CLI robusta e intuitiva.

### Funcionalidades:

- 📄 Parse de arquivos YAML e JSON
- 🔍 Validação de campos obrigatórios 
- 🖥️ Exibição completa da configuração
- 🖥️ Exibição filtrada apenas dos servidores
- ✅ Detecção automática do formato (YAML/JSON)
- 🧪 Testes unitários completos

### Estruturas suportadas:

- **Servidores**: nome, host, porta e réplicas
- **Banco de dados**: host, porta, usuário e senha

## Inicializando o Projeto

Criar o diretório estruturante para o CLI:

```bash
cobra-cli init
```

Adicionar subcomandos:

```bash
cobra-cli add parse
cobra-cli add server
```

## Executando os Comandos

### Comando Parse
Faz o parse completo do arquivo de configuração:

```bash
go run main.go parse --file example_config.yaml
```

### Comando Server
Exibe apenas a configuração dos servidores:

```bash
go run main.go server --file example_config.yaml
```

## Testes

Executar todos os testes:

```bash
go test ./...
```

Executar testes com detalhes:

```bash
go test ./config -v
```

## Makefile

```bash
make build    # Compila o binário
make test     # Executa os testes
make run      # Executa o programa
make clean    # Remove binários
```

## Docker

Construir a imagem:

```bash
make docker-build
```

Executar no container:

```bash
make docker-run
```

Ou manualmente:

```bash
docker build -t configparser .
docker run --rm -v $(PWD):/app configparser parse --file /app/example_config.yaml
```