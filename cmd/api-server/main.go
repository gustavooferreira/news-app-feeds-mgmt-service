package main

import (
	"fmt"
	"os"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/api"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/log"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/repository"
	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/lifecycle"
)

func main() {
	retCode := mainLogic()
	os.Exit(retCode)
}

func mainLogic() int {
	// Setup logger
	logger := core.NewAppLogger(os.Stdout, log.INFO)
	defer logger.Sync()

	logger.Info("APP starting")

	// Read config
	logger.Info("reading configuration", log.Field("type", "setup"))
	config := core.NewConfig()
	if err := config.LoadConfig(); err != nil {
		logger.Error(err.Error(), log.Field("type", "config"))
		return 1
	}

	// TODO: Set log level after reading config
	// something like this:
	// logger.SetLevel(config.Options.LogLevel)

	db, err := repository.NewDatabase()
	if err != nil {
		logger.Error(err.Error(), log.Field("type", "database"))
		return 1
	}

	server := api.NewServer(config.Webserver.Host, config.Webserver.Port, config.Options.DevMode, logger, db)

	// Spawn SIGINT listener
	go lifecycle.TerminateHandler(logger, server)

	// Listen for incoming requests -- app blocks here
	err = server.ListenAndServe()
	if err != nil {
		logger.Error(fmt.Sprintf("unexpected error while serving HTTP: %s", err))
		return 1
	}

	logger.Info("APP gracefully terminated")
	return 0
}
