package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/smartfactory-kl/idhub-key-tester/requester"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	buildInfo, _ := debug.ReadBuildInfo()

	logger := slog.New(handler)

	logger.Info(buildInfo.GoVersion)

	requester := requester.NewRequester(logger)

	requester.RunAll()
}
