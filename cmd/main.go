package main

import (
	"context"
	"flag"
	"time"

	"github.com/crxfoz/rater"
	"github.com/crxfoz/rater/providers/coinlayer"
	"github.com/crxfoz/rater/providers/coinmarketcap"
	"github.com/sirupsen/logrus"
)

func main() {
	port := flag.Int("port", 80, "Set port for web-server")
	period := flag.Int("period", 3600, "Period of parsing rates")

	flag.Parse()

	logger := logrus.WithField("service", "rater").WithField("module", "rater")

	// TODO: parse from env
	providers := []rater.Provider{
		coinlayer.New(""),
		coinmarketcap.New(""),
	}

	raterSvc, err := rater.New(
		providers,
		time.Second*time.Duration(*period),
		logger,
	)

	if err != nil {
		logger.WithError(err).Fatal("Rater couldn't be created")
		return
	}

	go func() {
		if err := raterSvc.Start(context.TODO(), nil); err != nil {
			logger.WithError(err).Error("rater stopped with error")
		}
	}()

	srv := rater.NewServer(raterSvc)
	if err := srv.Run(*port); err != nil {
		logger.WithError(err).Errorf("web-server stopped with error")
	}
}
