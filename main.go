package main

import (
	"log/slog"
	"os"

	"github.com/omaily/JWT/config"
)

var (
	logger *slog.Logger
)

func init() {
	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		// 	if a.Key == slog.SourceKey {
		// 		source := a.Value.Any().(*slog.Source)
		// 		source.File = filepath.Base(source.File)
		// 	}
		// 	return a
		// }}),
	)
	slog.SetDefault(logger)
}

func main() {
	conf := config.MustLoad()
	logger.Info("conf Server", struct {
		Addr string
		Port string
	}{
		Addr: conf.HTTPServer.Address,
		Port: conf.HTTPServer.Port,
	})
	logger.Debug("debug enables")

}
