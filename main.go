package main

import (
	"go-sample/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	endpointViper string = "/api/viper"
)

func main() {
	subLogger := k8sInfo()
	conf := config.LoadConfig(subLogger)

	router := gin.Default()

	router.GET(endpointViper, func(c *gin.Context) {
		c.JSON(200, conf)
		log.Trace().
			Str("endpoint", endpointViper).
			Msg("Método GET executado com sucesso.")
	})

	router.Run(":8080")
}

// busca e retorna informações do cluster/pod como um subnivel personalizado de log
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
