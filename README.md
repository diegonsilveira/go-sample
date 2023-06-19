## Projeto de testes utilizando o Viper, ZeroLog, GIN, Prometheus e K3d

#### Informações úteis

Neste sample foi utilizado a lib Viper (https://github.com/spf13/viper) para buscar configurações no arquivo config.yaml (raiz do projeto), a lib ZeroLog (https://github.com/rs/zerolog) para gerar logs em formato JSON, a lib GIN (https://github.com/gin-gonic/gin) para gerar um endpoint que retorna as informações contidas no arquivo de configuração (config.yaml) e o OpenTelemetry (https://github.com/open-telemetry/opentelemetry-go) em conjunto com o Prometheus para gerar métricas.

Além disso, o passo-a-passo a seguir ajuda a subir a imagem em um cluster local utilizando o K3D (https://k3d.io/v5.5.1/).

#### Como subir a aplicação local

Passo 1 - Comando para subir a aplicação:

    go run main.go

Passo 2 - Endpoint retornando as configurações:

    http://localhost:8080/api/viper

Passo 3 - Endpoint retornando as métricas da aplicação:

    http://localhost:8088/metrics

#### Como subir a aplicação no K3d

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

PASSO 9 - Remover cluster:

    k3d cluster delete my-cluster