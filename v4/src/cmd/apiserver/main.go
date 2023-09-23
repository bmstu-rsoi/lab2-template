package main

import (
	"log/slog"
	"os"

	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver"
	"github.com/migregal/bmstu-iu7-ds-lab2/apiserver/config"
)

func main() {
	lg := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.ReadConfig()
	if err != nil {
		lg.Error("[startup] failed to init config", "err", err.Error())
		os.Exit(1)
	}

	app, err := apiserver.New(lg, cfg)
	if err != nil {
		lg.Error("[startup] failed to init app", "err", err.Error())
		os.Exit(1)
	}

	app.Run(lg)
}
