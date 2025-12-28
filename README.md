# Platform Engineering CLI - Learning Project

Um projeto educacional focado em **estruturas básicas de CLI** em Go, demonstrando **boas práticas** para interfaces.

## Objetivo

Este repositório contém exercícios progressivos que ensinam como construir **ferramentas CLI robustas** para automação e monitoramento de infraestrutura, seguindo padrões da indústria.

## Tecnologias e Patterns

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

## Estrutura do Projeto

```
plataformrocks/
├── exerc<number>/      
└── README.md         
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


## Testes

Executar todos os testes:

```bash
go test ./...
```

Executar testes com detalhes:

```bash
go test ./config -v
```