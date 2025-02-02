package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pietdevries94/Kabisa/database"
	"github.com/pietdevries94/Kabisa/openapi"
	"github.com/pietdevries94/Kabisa/repositories"
	"github.com/pietdevries94/Kabisa/services"
	"github.com/rs/zerolog"

	// we autoload .env files into the environment variables, for convenience
	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	// The address the server will start listening
	listenAddress string
	// The timeout in seconds used when making http requests to external services
	httpClientTimeout string
	// The log level that will be shown in the console and stored in the log file (if enabled)
	logLevel string
	// The path to the log file. If this is not set, the application will only log to console
	logFilePath string
	// The connection string for sqlite
	sqliteDSN string
}

// application contains setup services, directly needed by it's httpHandler methods
type application struct {
	logger       *zerolog.Logger
	quoteService quoteService
}

func main() {
	config := initConfig()
	logger := initLogger(config)

	// init Application sets services, repositories and their dependencies
	app := initApplication(logger, config)

	srv, err := openapi.NewServer(app)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("failed to setup ogen api")
	}

	logger.Info().Str("address", config.listenAddress).Msg("starting server")
	err = http.ListenAndServe(config.listenAddress, srv)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("address", config.listenAddress).
			Msg("failed to start http server")
	}
}

// initConfig retrieves the config from the .env file or environment variables.
// All env variables are prefixed by KABISAQUOTE_ to prevent unwanted collision by other environment variables
//
// This function is very rudimentary and should be replaced if the number of configuration options grow or mutate often.
func initConfig() *config {
	// First, we initialize the config with default values
	conf := &config{
		listenAddress:     "127.0.0.1:3333",
		httpClientTimeout: "10",
		logLevel:          "info",
		logFilePath:       "",
		sqliteDSN:         "file::memory:?cache=shared",
	}

	// Next, we manually lookup the environment variables
	if val, found := os.LookupEnv("KABISAQUOTE_LISTEN_ADDRESS"); found {
		conf.listenAddress = val
	}
	if val, found := os.LookupEnv("KABISAQUOTE_HTTP_CLIENT_TIMEOUT"); found {
		conf.httpClientTimeout = val
	}
	if val, found := os.LookupEnv("KABISAQUOTE_LOG_LEVEL"); found {
		conf.logLevel = val
	}
	if val, found := os.LookupEnv("KABISAQUOTE_LOG_FILE_PATH"); found {
		conf.logFilePath = val
	}
	if val, found := os.LookupEnv("KABISAQUOTE_SQLITE_DSN"); found {
		conf.logFilePath = val
	}

	return conf
}

// initLogger creates a simple logger, which outputs to the console. The logLevel and logFilePath from the config are used.
// logLevel determines which minimum level of logs are shown in the console and written to the file.
// If logFilePath is set, the application also writes the logs as json to the path. Otherwise, only the console is used.
//
// If growing logs become a problem and an external log rotator is not a valid option, this function should rotate logs themselves.
// For now, it will only append logs to the given file path.
func initLogger(conf *config) *zerolog.Logger {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}

	// We create the basic logger as soon as possible, so we can use it for reporting errors later in the init
	logger := zerolog.New(consoleWriter).
		With().Timestamp().Caller().
		Logger()

	lvl, err := zerolog.ParseLevel(conf.logLevel)
	if err != nil {
		logger.Fatal().Err(err).Str("value", conf.logLevel).Msg("error when passing logLevel")
	}
	logger = logger.Level(lvl)

	// If the logFilePath is not set, we can safely early return
	if conf.logFilePath == "" {
		return &logger
	}

	logFile, err := os.OpenFile(
		conf.logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		logger.Fatal().Err(err).Str("path", conf.logFilePath).Msg("can't open/create given file path")
	}

	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	logger = logger.Output(multi)
	return &logger
}

// initApplication sets up the services, repositories and their dependencies
// It returns a struct which contains the logger and services to be used by it's httpHandler methods
func initApplication(logger *zerolog.Logger, conf *config) *application {
	db := database.Init(logger, conf.sqliteDSN)
	httpClient := initHttpClient(logger, conf)

	dummyJsonRepo := repositories.NewDummyJsonRepo(logger, httpClient)
	quoteGameRepo := repositories.NewQuoteGameRepo(logger, db)
	quoteService := services.NewQuoteService(logger, dummyJsonRepo, quoteGameRepo)

	return &application{
		logger:       logger,
		quoteService: quoteService,
	}
}

// initHttpClient creates a http client to be used for http requests to external services
// httpClientTimeout from the config is passed to the client
func initHttpClient(logger *zerolog.Logger, conf *config) *http.Client {
	timeoutInt, err := strconv.Atoi(conf.httpClientTimeout)
	if err != nil {
		logger.Fatal().Err(err).Str("value", conf.httpClientTimeout).Msg("could not parse set httpClientTimeout as int")
	}

	return &http.Client{
		Timeout: time.Duration(timeoutInt) * time.Second,
	}
}
