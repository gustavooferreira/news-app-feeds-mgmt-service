package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gustavooferreira/news-app-feeds-mgmt-service/pkg/core/log"
)

// NOTE: We would replace this with a proper config library like Viper

const AppPrefix = "NEWS_APP_FEEDS_MGMT"

// Configuration holds the entire configuration
type Configuration struct {
	Webserver WebserverConfiguration
	Options   OptionsConfiguration
}

// WebserverConfiguration holds configuration related to the webserver
type WebserverConfiguration struct {
	Host string
	Port int
}

// OptionsConfiguration holds general configuration
type OptionsConfiguration struct {
	// Development mode disables the panic recovery so we can see what was the actual problem.
	// and also, enables pprof
	DevMode bool

	LogLevel log.Level
}

// NewConfig returns new default configuration
func NewConfig() (config Configuration) {
	config.setDefaults()
	return config
}

// LoadConfig loads and validates config (from env vars)
func (config *Configuration) LoadConfig() (err error) {

	if webserverHost, ok := os.LookupEnv(AppPrefix + "_WEBSERVER_HOST"); ok {
		config.Webserver.Host = webserverHost
	}

	if webserverPort, ok := os.LookupEnv(AppPrefix + "_WEBSERVER_PORT"); ok {
		config.Webserver.Port, err = strconv.Atoi(webserverPort)
		if err != nil || config.Webserver.Port <= 0 || config.Webserver.Port > 1<<16-1 {
			return fmt.Errorf("configuration error: [webserver port] input not allowed <%s>", webserverPort)
		}
	}

	if devMode, ok := os.LookupEnv(AppPrefix + "_DEV_MODE"); ok {
		config.Options.DevMode, err = strconv.ParseBool(devMode)
		if err != nil {
			return fmt.Errorf("configuration error: [options devmode] unrecognizable boolean <%s>", devMode)
		}
	}

	if logLevel, ok := os.LookupEnv(AppPrefix + "_LOG_LEVEL"); ok {
		config.Options.LogLevel, err = ParseLogLevel(logLevel)
		if err != nil {
			return fmt.Errorf("configuration error: [options loglevel] unrecognized log level")
		}
	}

	return nil
}

func (config *Configuration) setDefaults() {
	// Webserver
	config.Webserver.Host = "127.0.0.1"
	config.Webserver.Port = 8080

	// Options
	config.Options.DevMode = false
	config.Options.LogLevel = log.INFO
}

func ParseLogLevel(level string) (logLevel log.Level, err error) {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		logLevel = log.DEBUG
	case "info":
		logLevel = log.INFO
	case "warning":
		logLevel = log.WARN
	case "error":
		logLevel = log.ERROR
	default:
		return 0, fmt.Errorf("log level unrecognised")
	}

	return logLevel, nil
}
