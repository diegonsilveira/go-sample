package main

import (
	"context"
	"go-sample/config"
	"go-sample/metrics"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	endpointViper string = "/api/viper"
)

func main() {
	// Configura o OpenTelemetry
	metrics.ConfigureOpentelemetry()
	ctx := context.Background()

	// Gera a métrica
	meter := otel.GetMeterProvider().Meter("request")
	requestsTotal, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total de requisições"),
	)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao gerar a métrica http_requests_total.")
	}

	subLogger := k8sInfo()

	// Carrega as configurações do arquivo config.yaml
	conf := config.LoadConfig(subLogger)

	router := gin.Default()

	router.GET(endpointViper, func(c *gin.Context) {
		c.JSON(200, conf)
		log.Trace().
			Str("endpoint", endpointViper).
			Msg("Método GET executado com sucesso.")

		// Incrementa a métrica de contagem de solicitações
		requestsTotal.Add(ctx, 1)
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
