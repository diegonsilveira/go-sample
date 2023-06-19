package main

import (
	"context"
	"go-sample/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	endpointViper string = "/api/viper"
)

func main() {
	ctx := context.Background()
	configureOpentelemetry()

	meter := otel.GetMeterProvider().Meter("request")
	requestsTotal, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total de requisições"),
	)
	if err != nil {
		panic(err)
	}

	subLogger := k8sInfo()

	//Carrega as configurações do arquivo config.yaml
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

// Expoe as métricas do Prometheus no endpoint localhost:8088/metrics utilizando o OpenTelemetry
func configureOpentelemetry() {
	router := gin.Default()

	if err := runtimemetrics.Start(); err != nil {
		panic(err)
	}

	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))

	otel.SetMeterProvider(provider)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	go func() {
		router.Run(":8088")
	}()
}
