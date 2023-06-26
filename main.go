package main

import (
	"context"
	"go-sample/api"
	"go-sample/config"
	_ "go-sample/docs"
	"go-sample/metrics"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	endpointViper string = "/api/viper"
)

// @title           Go Sample - CloudNative Team
// @version         1.0
// @description     This is a sample

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
func main() {
	gin.SetMode(gin.ReleaseMode)

	// Configura o OpenTelemetry
	metrics.ConfigureOpentelemetry()

	// Gera a métrica
	ctx := context.Background()
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
		api.GetConfigHandler(c, conf, requestsTotal, ctx, endpointViper)
	})

	// Gera o portal do Swagger (http://localhost:8080/swagger/index.html)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
