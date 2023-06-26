package api

import (
	"context"
	"go-sample/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/metric"
)

// GetConfigHandler retorna a configuração do endpoint.
// @Summary Retorna as configurações da aplicação.
// @Description Retorna as configurações da aplicação.
// @Tags sample
// @Accept json
// @Produce json
// @Success 200 {object} config.Configuration
// @Router /viper [get]
func GetConfigHandler(c *gin.Context, conf *config.Configuration, requestsTotal metric.Int64Counter, context context.Context, endpoint string) {
	c.JSON(http.StatusOK, conf)
	log.Trace().
		Str("endpoint", endpoint).
		Msg("Método GET executado com sucesso.")

	// Incrementa a métrica de contagem de solicitações
	requestsTotal.Add(context, 1)
}
