# Exercício 04 - CLI Kubernetes + Operator com Helm

Este exercício contém **duas funcionalidades principais**:

1. **CLI em Go** para comandos Kubernetes usando Cobra
2. **Kubernetes Operator** baseado em Helm Chart usando operator-sdk

## Estrutura do Projeto

```
exerc04/
├── cmd/                           # 🔧 CLI Kubernetes em Go (Cobra)
│   ├── k8s.go                     # Comando raiz 'k8s'
│   ├── list-pods.go               # Subcomando 'list' para pods
│   ├── root.go                    # Comando raiz principal
│   └── list-pods_test.go          # Testes do CLI
├── config/                        # ⚙️ Configurações Kubernetes (Kustomize)
│   ├── crd/                       # Custom Resource Definitions
│   ├── default/                   # Configuração padrão (namespace: exerc04-system)
│   ├── manager/                   # Deployment do operator
│   ├── rbac/                      # Roles e RoleBindings
│   └── samples/                   # Exemplos de Custom Resources
├── helm-charts/
│   └── visitors-helm/             # 📦 Helm Chart para o Operator
├── operator/                      # Chart original usado como base
├── watches.yaml                   # Configuração do helm-operator
├── Dockerfile                     # Container do operator
├── Makefile                       # Comandos de build/deploy
├── main.go                        # Código Go principal
└── go.mod                         # Dependências Go
```

---

## 🔧 PARTE 1: CLI Kubernetes em Go

### Funcionalidades do CLI

O CLI possui comandos para interagir com clusters Kubernetes:

- **`k8s list <namespace>`** - Lista pods em um namespace específico
- Funções auxiliares para listar namespaces
- Integração com kubeconfig (`~/.kube/config`)

### Uso do CLI

```bash
# Compilar o CLI
go build -o k8s-cli main.go

# Listar pods no namespace default
./k8s-cli k8s list default

# Listar pods no namespace kube-system
./k8s-cli k8s list kube-system
```

### Teste do CLI

```bash
# Executar testes
go test ./cmd -v

# Testar apenas o list-pods
go test ./cmd -run TestListPods -v
```

### Arquivos do CLI

#### `cmd/k8s.go`
Define o comando raiz `k8s` para operações do Kubernetes.

#### `cmd/list-pods.go`
Implementa:
- **`ListPods()`** - Lista pods em um namespace
- **`ListNamespaces()`** - Lista todos os namespaces
- Comando `list` que aceita namespace como argumento

#### `cmd/list-pods_test.go`
Testes abrangentes usando:
- **Fake Kubernetes client** para simular cluster
- **Testes de diferentes namespaces**
- **Validação de nomes de pods retornados**

---

## ⚙️ PARTE 2: Kubernetes Operator com Helm
│   ├── crd/                       # Custom Resource Definitions
│   ├── default/                   # Configuração padrão (namespace: exerc04-system)
│   ├── manager/                   # Deployment do operator
│   ├── rbac/                      # Roles e RoleBindings
│   └── samples/                   # Exemplos de Custom Resources
├── helm-charts/
│   └── visitors-helm/             # Helm Chart original
├── operator/                      # Chart original usado como base
├── watches.yaml                   # Configuração do helm-operator
├── local-watches.yaml            # Para execução local (removido)
├── Dockerfile                     # Container do operator
├── Makefile                       # Comandos de build/deploy
└── main.go                        # Código Go principal
```

## Passo a Passo da Criação

### 1. Iniciação do Operator baseado em Helm Chart existente

```bash
# Comando usado (na versão atual do operator-sdk):
operator-sdk init --plugins=helm --domain=example.com --group=example --version=v1 --kind=VisitorsApp --helm-chart=./operator
```

**Nota:** O comando antigo era `operator-sdk new` mas foi substituído por `operator-sdk init`.

### 2. Estrutura gerada automaticamente

O operator-sdk criou:
- **Kustomize structure** em `config/` (padrão moderno)
- **CRD** para `VisitorsApp` em `config/crd/bases/`
- **RBAC** completo em `config/rbac/`
- **Deployment** do operator em `config/manager/`
- **Samples** em `config/samples/example_v1_visitorsapp.yaml`
- **Dockerfile** configurado para helm-operator
- **Makefile** com targets para build e deploy

### 3. Aplicação dos CRDs no cluster

```bash
# Aplicar CRDs usando Kustomize
kubectl apply -k config/crd
```

### 4. Build e Deploy do Operator

#### 4.1 Construção da imagem Docker

```bash
# Build da imagem
docker build -t controller:latest .

# Load no cluster Kind (para desenvolvimento local)
kind load docker-image controller:latest
```

#### 4.2 Deploy no cluster

```bash
# Deploy completo usando Kustomize
kubectl apply -k config/default
```

**Importante:** O arquivo `config/manager/manager.yaml` foi modificado para incluir:
```yaml
imagePullPolicy: IfNotPresent  # Para usar imagem local no Kind
```

### 5. Verificação do Deployment

```bash
# Verificar se o operator está rodando
kubectl get pods -n exerc04-system

# Verificar logs
kubectl logs -n exerc04-system -l control-plane=controller-manager
```

### 6. Teste do Operator

```bash
# Aplicar um exemplo de Custom Resource
kubectl apply -f config/samples/example_v1_visitorsapp.yaml

# Verificar se os recursos foram criados pelo Helm
kubectl get all -l app.kubernetes.io/managed-by=Helm
```

## Arquivos Importantes

### `watches.yaml`
Configuração que define quais Custom Resources o operator deve monitorar:
```yaml
- group: example.example.com
  version: v1
  kind: VisitorsApp
  chart: helm-charts/visitors-helm  # Path relativo no container
```

### `config/samples/example_v1_visitorsapp.yaml`
Exemplo de Custom Resource que pode ser aplicado:
```yaml
apiVersion: example.example.com/v1
kind: VisitorsApp
metadata:
  name: visitorsapp-sample
spec:
  backend:
    size: 1
  frontend:
    title: Helm Installed Visitors Site
```

### `config/default/kustomization.yaml`
Configuração principal que define:
- **Namespace:** `exerc04-system`
- **Prefix:** `exerc04-`
- **Recursos:** CRD + RBAC + Manager + Metrics

## Comandos Úteis

### CLI Kubernetes
```bash
# Compilar e testar o CLI
go build -o k8s-cli main.go
./k8s-cli k8s list default
go test ./cmd -v
```

### Operator - Desenvolvimento Local
```bash
# Executar localmente (fora do cluster) - requer local-watches.yaml com paths absolutos
go run main.go

# Ou usando o binário compilado
./bin/helm-operator run --watches-file ./local-watches.yaml
```

### Deploy/Undeploy Operator
```bash
# Deploy tudo
kubectl apply -k config/default

# Remove tudo
kubectl delete -k config/default

# Remove apenas CRDs
kubectl delete -k config/crd
```

### Logs e Debug do Operator
```bash
# Logs do operator
kubectl logs -n exerc04-system deployment/exerc04-controller-manager

# Status dos pods
kubectl get pods -n exerc04-system

# Describe para debug
kubectl describe pod -n exerc04-system -l control-plane=controller-manager
```

## Problemas Encontrados e Soluções

### 1. **ImagePullBackOff**
**Problema:** Pod não conseguia fazer pull da imagem `controller:latest`
**Solução:** Build local + Kind load + `imagePullPolicy: IfNotPresent`

### 2. **Path absoluto vs relativo**
**Problema:** `watches.yaml` com paths diferentes para execução local vs container
**Solução:** 
- `watches.yaml`: path relativo para container
- `local-watches.yaml`: path absoluto para execução local

### 3. **Estrutura do Helm Chart**
**Problema:** Chart precisa estar na estrutura `helm-charts/visitors-helm/`
**Solução:** Reorganizar diretórios conforme esperado pelo operator

## Tecnologias Utilizadas

### CLI Kubernetes
- **Go** como linguagem principal
- **Cobra** para estrutura de comandos CLI
- **client-go** para interação com Kubernetes API
- **Testing** com fake Kubernetes client

### Operator
- **Operator SDK v1.x** (helm plugin)
- **Helm Operator** como runtime
- **Kustomize** para organização de manifests
- **Kind** para cluster local
- **Docker** para containerização

## Conceitos Aprendidos

### CLI Development
1. **Cobra framework** para CLIs em Go
2. **Kubernetes client-go** para API interactions  
3. **Testing with fakes** para simular clusters
4. **Kubeconfig management** para autenticação

### Operator Development

1. **Helm Operators** convertem Helm Charts em Kubernetes Operators
2. **Kustomize** é o padrão moderno para organização de manifests
3. **CRDs** definem novos tipos de recursos no Kubernetes
4. **RBAC** é gerado automaticamente baseado no Helm Chart
5. **OLM** é opcional - pode deployar com kubectl diretamente