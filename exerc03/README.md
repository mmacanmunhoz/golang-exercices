# Docker CLI Management Tool

Uma ferramenta CLI em Go para gerenciamento de containers Docker, implementando padrões de concorrência e integração com Docker API.

## Funcionalidades

### Container Management
- **list**: Lista todos os containers em execução
- **logs**: Visualiza logs de um container específico  
- **exec**: Executa comandos dentro de containers
- **health**: Verifica status de saúde de containers
- **stop**: Para todos os containers em execução
- **deploy**: Deploy de múltiplos containers via configuração YAML

### Deploy Concorrente
- Worker pool com 10 goroutines para processamento paralelo
- Pull de imagens simultâneo
- Criação e inicialização de containers em paralelo
- Configuração de redes e volumes

## Uso

### Comandos Básicos

```bash
# Compilar o projeto
go build -o docker-cli .

# Listar containers
./docker-cli container list

# Ver logs de um container
./docker-cli container logs <container-id>

# Executar comando em container
./docker-cli container exec <container-id> "ls -la"

# Verificar saúde do container
./docker-cli container health <container-id>

# Parar todos os containers
./docker-cli container stop
```

### Deploy de Containers

```bash
# Deploy via arquivo de configuração
./docker-cli container deploy --file deploy.yaml
```

### Estrutura do arquivo de configuração (deploy.yaml)

```yaml
servers:
  - name: "nginx-web-app"
    image: "nginx:1.21.4"
    ports:
      "80/tcp": "8080"
      "443/tcp": "8443"
    env:
      - "NGINX_HOST=localhost"
      - "NGINX_PORT=80"
    volumes:
      - "/host/nginx/html:/usr/share/nginx/html:ro"
      - "/host/nginx/conf:/etc/nginx/conf.d:ro"
    networks:
      - "bridge"
      - "web-network"

  - name: "redis-cache"
    image: "redis:7.0-alpine"
    ports:
      "6379/tcp": "6379"
    env:
      - "REDIS_PASSWORD=mypassword123"
    volumes:
      - "/host/redis/data:/data:rw"
    networks:
      - "bridge"
      - "backend-network"
```

## Tecnologias Utilizadas

- **Go 1.24+**: Linguagem principal
- **Cobra**: Framework para CLI
- **Docker API**: Integração com Docker daemon
- **YAML**: Parser de configuração
- **Worker Pool Pattern**: Concorrência eficiente
- **Context**: Timeout e cancelamento

## Padrões Implementados

### Concorrência
- Worker pool com canal de jobs
- Sync.WaitGroup para sincronização
- Context para timeout de operações

### Docker Integration
- Client API do Docker
- Container lifecycle management
- Port binding e network configuration
- Volume mounting

### Error Handling
- Tratamento robusto de erros
- Logs informativos
- Recuperação graceful de falhas

## Estrutura do Projeto

```
exerc03/
├── cmd/                    # Comandos CLI
│   ├── container.go       # Comando pai container
│   ├── container_exec.go  # Execução de comandos
│   ├── container_logs.go  # Visualização de logs
│   ├── deploy.go          # Deploy concorrente
│   ├── health.go          # Health checking
│   ├── list_containers.go # Listagem de containers
│   └── stop_containers.go # Parada de containers
├── config/                # Estruturas de configuração
│   └── types.go          # Tipos e structs
├── deploy.yaml           # Exemplo de configuração
├── main.go              # Entry point
└── README.md           # Este arquivo
```

## Build e Execução

```bash
# Instalar dependências
go mod tidy

# Compilar
go build -o docker-cli .

# Executar
./docker-cli container --help
```

## Dependências

- github.com/docker/docker: Docker API client
- github.com/docker/go-connections: Networking utilities
- github.com/spf13/cobra: CLI framework
- gopkg.in/yaml.v3: YAML parsing