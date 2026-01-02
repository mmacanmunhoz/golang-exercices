# Exercício 04 - CLI Kubernetes + Operator com Helm

Este exercício contém **duas funcionalidades principais**:

1. **CLI em Go** para comandos Kubernetes usando Cobra
2. **Kubernetes Operator** baseado em Helm Chart usando operator-sdk

## Estrutura do Projeto

```
exerc04/
├── cmd/                           # CLI Kubernetes em Go (Cobra)
│   ├── k8s.go                     # Comando raiz 'k8s'
│   ├── list-pods.go               # Subcomando 'list' para pods
│   ├── root.go                    # Comando raiz principal
│   └── list-pods_test.go          # Testes do CLI
├── config/                        # Configurações Kubernetes (Kustomize)
│   ├── crd/                       # Custom Resource Definitions
│   ├── default/                   # Configuração padrão (namespace: exerc04-system)
│   ├── manager/                   # Deployment do operator
│   ├── rbac/                      # Roles e RoleBindings
│   └── samples/                   # Exemplos de Custom Resources
├── helm-charts/
│   └── visitors-helm/             # Helm Chart para o Operator
├── operator/                      # Chart original usado como base
├── watches.yaml                   # Configuração do helm-operator
├── Dockerfile                     # Container do operator
├── Makefile                       # Comandos de build/deploy
├── main.go                        # Código Go principal
└── go.mod                         # Dependências Go
```

---

## PARTE 1: CLI Kubernetes em Go

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

## PARTE 2: Kubernetes Operator com Helm
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

---

## Versionamento de APIs (CRDs)

### Como adicionar novas versões ao CRD

Para **Helm Operators**, o versionamento é feito editando o CRD manualmente:

#### 1. Editar o CRD base
Arquivo: `config/crd/bases/example.example.com_visitorsapps.yaml`

```yaml
spec:
  group: example.example.com
  versions:
  # Versão v1 (existente)
  - name: v1
    served: true
    storage: true      # Versão de storage atual
    schema: # ... schema v1
    
  # Nova versão v2 
  - name: v2
    served: true
    storage: false     # Ainda não é storage
    schema:
      openAPIV3Schema:
        properties:
          spec:
            properties:
              backend:
                properties:
                  size:
                    type: integer
                    minimum: 1      # Validações aprimoradas
                    maximum: 10
              monitoring:           # Novo campo em v2
                properties:
                  enabled:
                    type: boolean
```

#### 2. Aplicar CRD atualizado
```bash
kubectl apply -k config/crd
```

#### 3. Testar múltiplas versões

**CR v1 (formato original):**
```yaml
apiVersion: example.example.com/v1
kind: VisitorsApp
metadata:
  name: visitorsapp-v1-sample
spec:
  backend:
    size: 1
  frontend:
    title: "Site v1"
```

**CR v2 (formato aprimorado):**
```yaml
apiVersion: example.example.com/v2
kind: VisitorsApp
metadata:
  name: visitorsapp-v2-sample
spec:
  backend:
    size: 3
    version: "v1.2.0"
  frontend:
    title: "Enhanced Site v2"
    replicas: 2
  monitoring:            # Campo disponível apenas em v2
    enabled: true
    metrics: true
```

#### 4. Comandos de teste

```bash
# Listar CRs de versões específicas
kubectl get visitorsapps.v1.example.example.com
kubectl get visitorsapps.v2.example.example.com

# Verificar versões disponíveis
kubectl api-resources | grep visitors

# Aplicar CRs de diferentes versões
kubectl apply -f config/samples/example_v1_visitorsapp.yaml
kubectl apply -f - <<EOF
apiVersion: example.example.com/v2
kind: VisitorsApp
metadata:
  name: test-v2
spec:
  backend:
    size: 5
  frontend:
    title: "Test v2"
  monitoring:
    enabled: true
EOF
```

### Estratégias de Migration

#### Deprecação gradual:
1. **Adicionar v2** mantendo v1 como storage
2. **Testar v2** em paralelo com v1
3. **Migrar CRs** gradualmente para v2
4. **Marcar v1 como deprecated**
5. **Trocar storage** para v2
6. **Remover v1** em release futuro

#### Conversion (opcional):
Para conversão automática entre versões, seria necessário implementar **conversion webhooks** - mais complexo para Helm operators.

### Validações Avançadas em v2

O schema v2 inclui validações aprimoradas:

```yaml
properties:
  backend:
    properties:
      size:
        minimum: 1        # Tamanho mínimo
        maximum: 10       # Tamanho máximo
      version:
        pattern: '^v[0-9]+\.[0-9]+\.[0-9]+$'  # Formato semver
  frontend:
    properties:
      title:
        minLength: 1      # Não pode ser vazio
        maxLength: 100    # Limite de caracteres
required:
- backend               # Campos obrigatórios
- frontend
```

**Resultado**: Ambas versões funcionam simultaneamente no mesmo operator!

---

## Estratégias de Versionamento de Operators

### Problema: Chart Embutido na Imagem

Como o Helm Chart é **copiado inteiro** para dentro da imagem Docker do operator via Dockerfile, cada evolução do chart requer uma nova imagem do operator.

### Estratégia Recomendada: Múltiplas Imagens (MAIS PRÁTICA)

**RECOMENDAÇÃO:** Esta abordagem é significativamente mais prática e segura para produção.

Manter **imagens separadas** para cada versão major do operator/chart:

```
mmacanmunhoz/controller-demo:v1.0  # Chart v1 + CRD v1 only  
mmacanmunhoz/controller-demo:v2.0  # Chart v2 + CRD v2 focus
```

**Por que é mais prática:**
- Isolamento total entre versões (sem riscos de breaking changes)
- Deploy e rollback simplificados (apenas troca de imagem)
- Charts limpos sem código legacy
- Facilita testes A/B e deployments side-by-side

### Implementação Prática

#### 1. Preparar Versão v1.0 (Baseline)

```bash
# Guardar versão v1 estável
git tag v1.0-chart
git push origin v1.0-chart

# Build imagem final v1.0
docker build -t mmacanmunhoz/controller-demo:v1.0 .
docker push mmacanmunhoz/controller-demo:v1.0
```

#### 2. Evoluir para v2.0 (Breaking Changes)

**a) Limpar values.yaml para v2:**
```yaml
# helm-charts/visitors-helm/values.yaml
backend:
  size: 3               # Valores otimizados para v2
  version: "v1.2.0"     # Novo campo

frontend:
  title: "Enhanced Visitors Site v2"
  replicas: 2           # Novo campo

monitoring:             # Funcionalidade nova
  enabled: true
  metrics: true
```

**b) Atualizar CRD - v2 como storage:**
```yaml
# config/crd/bases/example.example.com_visitorsapps.yaml
versions:
- name: v1
  served: true
  storage: false     # v1 deprecated
  deprecated: true
  
- name: v2
  served: true  
  storage: true      # v2 é storage version
  schema: # ... schema aprimorado com validações
```

**c) Build nova imagem:**
```bash
docker build -t mmacanmunhoz/controller-demo:v2.0 .
docker push mmacanmunhoz/controller-demo:v2.0
```

### Migration Strategy

**NOTA:** Para produção, a estratégia de múltiplas imagens é mais segura e prática que tentar manter retrocompatibilidade complexa.

#### Opção 1: Side-by-side (Recomendado)

```bash
# Cluster A: Continua com v1.0
kubectl patch deployment exerc04-controller-manager \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","image":"mmacanmunhoz/controller-demo:v1.0"}]}}}}'

# Cluster B: Migra para v2.0  
kubectl patch deployment exerc04-controller-manager \
  -p '{"spec":{"template":{"spec":{"containers":[{"name":"manager","image":"mmacanmunhoz/controller-demo:v2.0"}]}}}}'
```

#### Opção 2: In-place Migration

```bash
# 1. Deploy imagem v2.0 (com CRD suportando v1+v2)
kubectl apply -k config/default

# 2. Migrar CRs existentes v1 → v2
kubectl get visitorsapps.v1.example.example.com -o yaml | \
  sed 's/apiVersion: example.example.com\/v1/apiVersion: example.example.com\/v2/' | \
  kubectl apply -f -

# 3. Validar funcionamento
kubectl get visitorsapps.v2.example.example.com

# 4. Limpar CRs v1 antigos (opcional)
kubectl delete visitorsapps.v1.example.example.com --all
```

### Vantagens desta Estratégia

#### **Charts Limpos**
- Cada imagem foca em sua versão específica
- Sem código "morto" ou campos legacy desnecessários  
- Values.yaml otimizado para cada versão

#### **Versionamento Explícito**  
- Tag da imagem = versão do chart + operator
- Sem ambiguidade sobre qual versão está rodando
- Facilita troubleshooting e auditoria

#### **Rollback Simplificado**
- Rollback = troca de imagem Docker
- Não precisa reverter CRDs ou configurações
- Teste de versões em paralelo

#### **CI/CD Otimizado**
```yaml
# .github/workflows/operator.yml
strategy:
  matrix:
    version: [v1.0, v2.0]
    
steps:
- name: Build Operator ${{ matrix.version }}
  run: |
    git checkout ${{ matrix.version }}-chart
    docker build -t mmacanmunhoz/controller-demo:${{ matrix.version }} .
```

### Exemplo de Uso

#### Deploy v1.0 (Produção Estável)
```yaml
apiVersion: example.example.com/v1
kind: VisitorsApp
metadata:
  name: prod-app
spec:
  backend:
    size: 2
  frontend:
    title: "Production Site"
```

#### Deploy v2.0 (Features Avançadas)  
```yaml
apiVersion: example.example.com/v2
kind: VisitorsApp
metadata:
  name: enhanced-app
spec:
  backend:
    size: 5
    version: "v2.1.0"
  frontend:
    title: "Enhanced Production Site"
    replicas: 3
  monitoring:
    enabled: true
    metrics: true
```

### Considerações

#### **Manutenção**
- Manter duas branches/imagens requer mais overhead
- Patches de segurança podem precisar ser aplicados em ambas

#### **Transição**
- Período de grace para migration v1 → v2
- Comunicação clara sobre deprecation timeline
- Documentação de breaking changes

**Esta estratégia garante evolução controlada e deployment confiável dos operators em produção.**

---

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

### API Versioning
1. **Múltiplas versões** podem coexistir no mesmo CRD
2. **Validações OpenAPI** podem ser diferentes entre versões
3. **Storage version** define qual versão é armazenada no etcd
4. **Served versions** definem quais versões aceitam requests
5. **Migration gradual** permite evolução sem breaking changes
6. **Backward compatibility** mantém CRs antigos funcionando