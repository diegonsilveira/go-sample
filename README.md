## Projeto de testes utilizando o Viper, ZeroLog, GIN, OpenTelemetry Swaggo e K3d

### Informações úteis

Neste sample foram utilizadas algumas bibliotecas para desenvolvimento de experimentações. A seguir a lista de libs:

#### Viper

Busca de configurações no arquivo config.yaml e possibilidade de carregamento por variáveis de ambiente.

> Documentação: https://github.com/spf13/viper

#### ZeroLog

Geração de logs em formato JSON.

> Documentação: https://github.com/rs/zerolog

#### GIN

Geração do endpoint para testes que retorna as informações contidas no arquivo de configuração.

> Documentação: https://github.com/gin-gonic/gin

#### Swaggo/Swag

Geração da documentação da API de testes.

> Documentação: https://github.com/swaggo/swag

#### OpenTelemetry e Prometheus

Geração e publicação de métricas.

> Documentações: https://github.com/open-telemetry/opentelemetry-go

#### K3D

Utilizado para gerar um cluster local para testes.

> Documentação: https://k3d.io/v5.5.1/

### Como subir a aplicação local

Passo 1 - Comando para subir a aplicação:

    go run main.go

Passo 2 - Endpoint retornando as configurações:

    http://localhost:8080/api/viper

Passo 3 - Endpoint retornando as métricas da aplicação:

    http://localhost:8088/metrics

Passo 4 - Acesso a documentação das APIs (Swagger)

    http://localhost:8080/swagger/index.html

> Para atualizar as documentações, executar o comando "swag init".

### Como subir a aplicação no K3d

PASSO 1 - Comando para buildar o projeto:

    go build -o sample

Execução:

    ./sample

PASSO 2 - Construir imagem docker:

    docker build -t {{endereço do registro}}/sample:v1 .

PASSO 3 - Enviar imagem para o hub.docker.com:

    docker push {{endereço do registro}}/sample:v1

PASSO 4 - Criar cluster Kubernetes (K3d):

    k3d cluster create my-cluster

PASSO 5 - Implantar aplicação no cluster:

* Sem variáveis de ambiente:

    kubectl create deployment sample --image={{endereço do registro}}/sample:v1

* Com variáveis de ambiente definidas no arquivo substituindo os valores do config.yaml:

    kubectl apply -f deployment.yaml

PASSO 6 - Expor serviço:

    kubectl expose deployment sample --port=8080 --target-port=8080 --type=LoadBalancer

PASSO 7 - Verificar IP da aplicação:

    kubectl get service sample

PASSO 8 - Para acessar a aplicação (external-ip):

    http://{{external-ip}}:8080/api/viper
    http://{{external-ip}}:8080/metrics
    http://{{external-ip}}:8080/swagger/index.html

PASSO 9 - Remover cluster:

    k3d cluster delete my-cluster