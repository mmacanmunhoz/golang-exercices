# Platform Rocks - Kubernetes & Go Learning Projects

Este repositório contém uma série de exercícios práticos para aprender Kubernetes, Go, Docker e desenvolvimento de Operators.

## 🚀 Projetos Incluídos

Este repositório contém exercícios progressivos que ensinam como construir **ferramentas CLI robustas** para automação e monitoramento de infraestrutura, seguindo padrões da indústria.

## Tecnologias e Patterns

### **Core Technologies**
- **Go 1.24+** - Performance e concorrência nativa
- **Cobra CLI** - Framework robusto para interfaces de linha de comando  
- **Docker API** - Integração com containers e orquestração
- **YAML/JSON** - Configuração estruturada e saída de dados

### **Concurrency Patterns**
- **Worker Pool** - Processamento paralelo eficiente
- **Fan-out/Fan-in** - Distribuição e agregação de trabalho
- **Context** - Timeout e cancelamento de operações
- **Channels** - Comunicação segura entre goroutines

### **Platform Engineering Practices**
- **Health Checking** - Monitoramento de serviços e containers
- **Configuration Management** - Parsing e validação de configs YAML/JSON
- **Container Orchestration** - Deploy e gerenciamento de containers
- **Structured Logging** - Saídas formatadas para integração
- **Error Handling** - Tratamento robusto de falhas
- **Performance Monitoring** - Medição de response time e health status

## Estrutura do Projeto

```
plataformrocks/
├── exerc01/           # Config Parser CLI + Docker
├── exerc02/           # Health Checker + Performance Monitor  
├── exerc03/           # Docker CLI Management Tool
└── README.md          # Documentação geral
```

## Exercícios

### **Exercício 01: Config Parser CLI**
**Foco:** Fundamentos CLI + Parsing + Docker

```bash
cd exerc01/configparser
go run main.go parse --file example_config.yaml
go run main.go server --file example_config.yaml
```

**Conceitos:**
- CLI com subcomandos
- Parsing YAML/JSON
- Validação de configuração
- Containerização

---

### **Exercício 02: Health Checker + Speed Monitor**
**Foco:** Concorrência + HTTP + Monitoring

```bash
cd exerc02
go run main.go health --file example_config.yaml   
go run main.go response --file example_config.yaml  
```

**Conceitos:**
- Worker pool concorrente (10 workers)
- HTTP health checking
- Performance monitoring
- JSON structured output
- Context timeout (5s)

---

### **Exercício 03: Docker CLI Management**
**Foco:** Docker API + Concorrência + Container Orchestration

```bash
cd exerc03
go build -o docker-cli .
./docker-cli container list
./docker-cli container logs <container-id>
./docker-cli container exec <container-id> <command>
./docker-cli container deploy --file deploy.yaml
./docker-cli container health <container-id>
./docker-cli container stop
```

**Conceitos:**
- Docker API integration
- Worker pool para deployment paralelo
- Container lifecycle management
- YAML configuration parsing
- Health monitoring
- Command execution em containers
- Port binding e network configuration


## Testes

Executar todos os testes:

```bash
go test ./...
```

Executar testes com detalhes:

```bash
go test ./config -v
```