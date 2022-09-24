package main

import (
	"log"

	"github.com/vzemtsov/config-reloader/config"
	"github.com/vzemtsov/config-reloader/internal/app"
	"github.com/vzemtsov/config-reloader/pkg/metrics"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Run(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if !*cfg.InitMode {
		err = metrics.Run()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
