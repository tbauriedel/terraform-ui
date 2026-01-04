package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/tbauriedel/resource-nexus-core/internal/app"
	"github.com/tbauriedel/resource-nexus-core/internal/authentication"
)

func main() {
	_ = os.Setenv("LANG", "C")

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	var (
		adminPassword string
		configPath    string
	)

	flag.StringVar(&configPath, "config", "config.json", "Config file")
	flag.StringVar(&adminPassword, "admin-password", "", "Admin password")
	flag.Parse()

	conf, err := app.LoadConfig(configPath)
	if err != nil {
		slog.Error(err.Error())
		app.Exit(nil, 1)
	}

	// stdout is hardcoded here and level to INFO
	conf.Logging.Type = "stdout"
	conf.Logging.Level = "info"

	// we do not need a logfile here because stdout is hardcoded.
	db, logger, _, err := app.Bootstrap(conf)
	if err != nil {
		slog.Error(err.Error())
		app.Exit(nil, 1)
	}

	// close database connection on exit of main
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create admin user
	err = authentication.CreateAdminUser(adminPassword, conf.Security.PasswordHashing, db, ctx)
	if err != nil {
		logger.Error(err.Error())
		app.Exit(nil, 1)
	}

	logger.Info("admin user created. username: admin")
}
