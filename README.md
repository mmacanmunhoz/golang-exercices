# Platform Engineering CLI - Learning Project

Um projeto educacional focado em **estruturas básicas de CLI** em Go, demonstrando **boas práticas** para interfaces de **engenharia de plataforma**.

## 🎯 Objetivo

Este repositório contém exercícios progressivos que ensinam como construir **ferramentas CLI robustas** para automação e monitoramento de infraestrutura, seguindo padrões da indústria.

## 🚀 Tecnologias e Patterns

### **Core Technologies**
- **Go 1.24+** - Performance e concorrência nativa
- **Cobra CLI** - Framework robusto para interfaces de linha de comando
- **YAML/JSON** - Configuração estruturada e saída de dados

### **Concurrency Patterns**
- **Worker Pool** - Processamento paralelo eficiente
- **Fan-out/Fan-in** - Distribuição e agregação de trabalho
- **Context** - Timeout e cancelamento de operações
- **Channels** - Comunicação segura entre goroutines

### **Platform Engineering Practices**
- **Health Checking** - Monitoramento de serviços
- **Configuration Management** - Parsing e validação de configs
- **Structured Logging** - Saídas JSON para integração
- **Error Handling** - Tratamento robusto de falhas
- **Performance Monitoring** - Medição de response time

## 📁 Estrutura do Projeto

```
plataformrocks/
├── exerc01/           # Básico: Config Parser + Docker
├── exerc02/           # Avançado: Health Checker + Website Speed Monitor
├── exerc03/           # Concorrência: Worker Pool Pattern
└── README.md         # Este arquivo
```

## 🏋️ Exercícios

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
go run main.go health --file example_config.yaml    # Health checking
go run main.go parse --file example_config.yaml     # Config parsing
```

**Conceitos:**
- Worker pool concorrente (10 workers)
- HTTP health checking
- Performance monitoring
- JSON structured output
- Context timeout (5s)

---

### **Exercício 03: Worker Pool Pattern**
**Foco:** Concorrência pura + Patterns

```bash
cd exerc03
go run main.go    # Demonstra worker pool com 3 workers
```

**Conceitos:**
- Goroutines + WaitGroup
- Channel communication
- Graceful shutdown

## 🛠️ Como Usar

### **Pré-requisitos**
```bash
# Instalar Go 1.24+
go version

# Clonar projeto
git clone <repo-url>
cd plataformrocks
```

### **Execução**
```bash
# Testar cada exercício
cd exerc01/configparser && go run main.go --help
cd exerc02 && go run main.go --help
cd exerc03 && go run main.go
```

### **Build e Deploy**
```bash
# Docker (Exercício 01)
cd exerc01
docker build -t configparser .
docker run --rm -v $(PWD):/app configparser parse --file /app/example_config.yaml

# Binário nativo
go build -o bin/health-checker ./exerc02/main.go
./bin/health-checker health --file exerc02/example_config.yaml
```

## 📊 Features Implementadas

### **Exercício 02 - Funcionalidades Completas**

#### **Health Checking**
- ✅ **Concorrência** - 10 workers simultâneos
- ✅ **HTTP Monitoring** - GET requests com timeout
- ✅ **JSON Output** - Structured logging
- ✅ **Dynamic Fields** - Status, Worker ID, Timestamp
- ✅ **Error Handling** - Network failures, timeouts

#### **Website Speed Monitoring**
- ✅ **Performance Timing** - Response time measurement
- ✅ **Threshold Checking** - Fast/slow classification
- ✅ **Worker Pool** - Parallel execution
- ✅ **Context Control** - Timeout management

## 🎓 Conceitos de Platform Engineering

### **Reliability**
- Health checking automático
- Timeout e retry patterns
- Error handling robusto

### **Observability**
- Structured JSON logging
- Performance metrics
- Worker tracking

### **Scalability**
- Concorrência nativa
- Resource pooling
- Efficient batching

### **Developer Experience**
- CLI intuitiva com `--help`
- Clear error messages
- Flexible configuration

## 📈 Progressão de Aprendizado

1. **Básico** → CLI estruturada + Config management
2. **Intermediário** → HTTP monitoring + JSON output  
3. **Avançado** → Concorrência + Performance + Patterns

## 🔧 Extensões Futuras

- [ ] Fan-in pattern implementation
- [ ] Prometheus metrics export
- [ ] gRPC health checking
- [ ] Kubernetes service discovery
- [ ] Configuration hot-reload
- [ ] Rate limiting
- [ ] Circuit breaker pattern

## 📚 Aprendizado

Este projeto demonstra como construir **ferramentas CLI profissionais** para **Platform Engineering**, cobrindo desde parsing básico até **monitoramento concorrente** em escala.

**Ideal para:** DevOps Engineers, Platform Engineers, SREs aprendendo Go.

---

**Desenvolvido como projeto educacional para demonstrar boas práticas em CLI tools e concorrência em Go.**