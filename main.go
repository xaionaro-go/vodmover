package main

import (
	"context"
	"io"
	"os"

	"github.com/facebookincubator/go-belt/tool/logger"
	xlogrus "github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/spf13/pflag"
	"github.com/xaionaro-go/vodmover/pkg/vodmover"
	"github.com/xaionaro-go/vodmover/pkg/xpath"
	"gopkg.in/yaml.v3"
)

func main() {
	logLevel := logger.LevelInfo
	pflag.Var(&logLevel, "log-level", "")
	logFile := pflag.String("log-file", "~/vodmover.log", "")
	configPath := pflag.String("config-path", "~/vodmover.yaml", "")
	pflag.Parse()

	l := xlogrus.DefaultLogrusLogger()

	if *logFile != "" {
		logFileExpanded, err := xpath.Expand(*logFile)
		if err != nil {
			l.Errorf("unable to expand path '%s': %v", logFileExpanded, err)
			return
		}
		f, err := os.OpenFile(logFileExpanded, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0750)
		if err != nil {
			l.Errorf("failed to open log file '%s': %v", logFileExpanded, err)
			return
		}
		l.SetOutput(io.MultiWriter(os.Stderr, f))
	}

	ctx := logger.CtxWithLogger(context.Background(), xlogrus.New(l).WithLevel(logLevel))

	configPathExpanded, err := xpath.Expand(*configPath)
	if err != nil {
		logger.Fatalf(ctx, "unable to expand the config path: %v", err)
	}

	configBytes, err := os.ReadFile(configPathExpanded)
	if err != nil {
		logger.Fatalf(ctx, "unable to read the config: %v", err)
	}

	cfg := vodmover.Config{}
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		logger.Fatalf(ctx, "unable to parse the config: %v", err)
	}

	mover := vodmover.New(cfg)

	err = mover.Serve(ctx)
	logger.Fatalf(ctx, "unable to run the VOD mover service: %v", err)
}
