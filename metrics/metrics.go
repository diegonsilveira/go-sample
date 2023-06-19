package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Expoe as métricas do Prometheus no endpoint localhost:8088/metrics utilizando o OpenTelemetry
func ConfigureOpentelemetry() {
	router := gin.Default()

	if err := runtimemetrics.Start(); err != nil {
		log.Error().Err(err).Msg("Erro ao iniciar o OpenTelemetry.")
	}

	exporter, err := prometheus.New()
	if err != nil {
		log.Error().Err(err).Msg("Erro ao criar as configurações do Prometheus.")
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))

	otel.SetMeterProvider(provider)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	go func() {
		router.Run(":8088")
	}()
}
