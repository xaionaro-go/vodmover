package main

import (
	"context"
	"os"

	"github.com/facebookincubator/go-belt/tool/logger"
	xlogrus "github.com/facebookincubator/go-belt/tool/logger/implementation/logrus"
	"github.com/spf13/pflag"
	"github.com/xaionaro-go/vodmover/pkg/vodmover"
	"gopkg.in/yaml.v3"
)

func main() {
	logLevel := logger.LevelInfo
	pflag.Var(&logLevel, "log-level", "")
	configPath := pflag.String("config-path", "~/.vodmover.yaml", "")
	pflag.Parse()
	ctx := logger.CtxWithLogger(context.Background(), xlogrus.Default().WithLevel(logLevel))

	configBytes, err := os.ReadFile(*configPath)
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
