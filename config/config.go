package config

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Configuration struct {
	Environment string

	Server struct {
		Host string
		Port int
	}

	Log struct {
		Level int
	}
}

const (
	configFile string = "config.yaml"
)

func LoadConfig(logger zerolog.Logger) *Configuration {
	var config *Configuration

	conf := viper.GetViper()

	viper.SetConfigFile(configFile) // arquivo de configurações
	viper.AddConfigPath(".")        // caminho para o arquivo de configuração

	logger.Trace().
		Str("configFile", configFile).
		Msg("Informações do arquivo de configurações.")

	//sobreescreve os valores do arquivo com variáveis de ambiente (ENV)
	viper.AutomaticEnv()

	//altera as variáveis de ambiente (ENV) trocando do padrão do SO (_) para o padrão do viper (.)
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	err := conf.ReadInConfig()
	if err != nil {
		log.Panic().Err(err).Msg("Erro ao ler arquivo de configuração.")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Error().Err(err).Msg("Erro ao decodificar arquivo de configuração.")
	}

	// Verifica se a propriedade log foi definida
	if !viper.IsSet("log.level") {
		// Define o valor INFO como padrão
		viper.SetDefault("log.level", 1)
	}

	//seta o level global do log de acordo com a propriedade definida no yaml
	zerolog.SetGlobalLevel(zerolog.Level(config.Log.Level))

	logger.Trace().
		Str("environment", config.Environment).
		Str("server.host", config.Server.Host).
		Int("server.port", config.Server.Port).
		Int("log.level", config.Log.Level).
		Msg("Configurações carregadas.")

	return config
}
