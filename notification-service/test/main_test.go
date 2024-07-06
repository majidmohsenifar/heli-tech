package test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/majidmohsenifar/heli-tech/notification-service/config"
	"github.com/majidmohsenifar/heli-tech/notification-service/logger"
	"github.com/spf13/viper"
)

var loggerService *slog.Logger

var viperConfig *viper.Viper

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)

}

func getLogger() *slog.Logger {
	if loggerService != nil {
		return loggerService
	}
	loggerService = logger.NewLogger()
	return loggerService
}

func getViperConfig() *viper.Viper {
	if viperConfig != nil {
		return viperConfig
	}
	viperConfig = config.NewViper("../config/")
	return viperConfig
}
