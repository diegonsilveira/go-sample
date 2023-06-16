package main

import (
	"go-sample/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	endpointViper string = "/api/viper"
)

var (
	requestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de requisições",
	})
)

func main() {
	//Registra a métrica criada
	prometheus.MustRegister(requestsTotal)

	subLogger := k8sInfo()

	//Carrega as configurações do arquivo config.yaml
	conf := config.LoadConfig(subLogger)

	router := gin.Default()

	//Expoe as métricas do Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET(endpointViper, func(c *gin.Context) {
		c.JSON(200, conf)
		log.Trace().
			Str("endpoint", endpointViper).
			Msg("Método GET executado com sucesso.")

		// Incrementa a métrica de contagem de solicitações
		requestsTotal.Inc()
	})

	router.Run(":8080")
}

// Busca e retorna informações do cluster/pod como um subnivel personalizado de log
func k8sInfo() zerolog.Logger {
	clusterName := os.Getenv("K3D_CLUSTER_NAME")
	clusterAddress := os.Getenv("K3D_CLUSTER_HOST")
	clusterNodes := os.Getenv("K3D_CLUSTER_NODES")
	podName := os.Getenv("HOSTNAME")

	logger := log.With().
		Str("K3D_CLUSTER_NAME", clusterName).
		Str("K3D_CLUSTER_HOST", clusterAddress).
		Str("K3D_CLUSTER_NODES", clusterNodes).
		Str("HOSTNAME", podName).Logger()

	return logger
}
