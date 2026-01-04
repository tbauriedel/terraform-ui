package app

import (
	"log/slog"
	"os"

	"github.com/tbauriedel/resource-nexus-core/internal/common/fileutils"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/database"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

// LoadConfig loads the configuration from the given path.
func LoadConfig(configPath string) (config.Config, error) {
	var (
		conf config.Config
		err  error
	)

	// Load server config
	if !fileutils.FileExists(configPath) {
		slog.Info("no config file found or provided. loading default configuration")

		conf = config.LoadDefaults()
	} else {
		slog.Info("loading config from file", "file", configPath)

		conf, err = config.LoadFromJSONFile(configPath)
		if err != nil {
			return config.Config{}, err //nolint:wrapcheck
		}
	}

	slog.Debug("starting with configuration", "config", conf.GetConfigRedacted())

	return conf, nil
}

// Bootstrap initializes the application.
//
// Returns the database connection, and logger.
// Returns an error if something went wrong.
func Bootstrap(conf config.Config) (database.Database, *logging.Logger, *os.File, error) { //nolint:ireturn
	var (
		err     error
		db      database.Database
		logger  *logging.Logger
		logfile *os.File
	)

	// init logger
	if conf.Logging.Type == "file" {
		slog.Debug("creating logger for type 'file'")

		// open logfile
		logfile, err = fileutils.OpenFile(conf.Logging.File)
		if err != nil {
			return nil, nil, logfile, err //nolint:wrapcheck
		}

		defer func(logfile *os.File) {
			_ = logfile.Close()
		}(logfile)

		// create logger for type 'file'
		logger = logging.NewLoggerFile(conf.Logging, logfile)
	} else {
		slog.Debug("creating logger for type 'stdout'")

		// create logger for type 'stdout'
		logger = logging.NewLoggerStdout(conf.Logging)
	}

	logger.Info("initializing database connection")

	// create the database connection
	db, err = database.NewDatabase(conf.Database, logger)
	if err != nil {
		return nil, nil, logfile, err //nolint:wrapcheck
	}

	// test database connection
	err = db.TestConnection()
	if err != nil {
		return nil, nil, logfile, err //nolint:wrapcheck
	}

	logger.Info("database connection established and tested successfully")

	return db, logger, logfile, nil
}

// Exit closes the logfile and exits the application with the given code.
func Exit(logfile *os.File, code int) {
	if logfile != nil {
		_ = logfile.Close()
	}

	os.Exit(code)
}
