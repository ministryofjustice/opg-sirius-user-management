package main

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-go-common/env"
	"github.com/ministryofjustice/opg-go-common/telemetry"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/server"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	ctx := context.Background()
	logger := telemetry.NewLogger("opg-sirius-user-management")

	if err := run(ctx, logger); err != nil {
		logger.Error("fatal startup error", slog.Any("err", err.Error()))
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	port := getEnv("PORT", "8080")
	webDir := getEnv("WEB_DIR", "web")
	siriusURL := getEnv("SIRIUS_URL", "http://localhost:9001")
	siriusPublicURL := getEnv("SIRIUS_PUBLIC_URL", "")
	prefix := getEnv("PREFIX", "")
	exportTraces := env.Get("TRACING_ENABLED", "0") == "1"

	layouts, _ := template.
		New("").
		Funcs(map[string]interface{}{
			"join": func(sep string, items []string) string {
				return strings.Join(items, sep)
			},
			"contains": func(xs []string, needle string) bool {
				for _, x := range xs {
					if x == needle {
						return true
					}
				}

				return false
			},
			"prefix": func(s string) string {
				return prefix + s
			},
			"sirius": func(s string) string {
				return siriusPublicURL + s
			},
		}).
		ParseGlob(webDir + "/template/layout/*.gotmpl")

	files, _ := filepath.Glob(webDir + "/template/*.gotmpl")
	tmpls := map[string]*template.Template{}

	for _, file := range files {
		tmpls[filepath.Base(file)] = template.Must(template.Must(layouts.Clone()).ParseFiles(file))
	}

	shutdown, err := telemetry.StartTracerProvider(ctx, logger, exportTraces)
	defer shutdown()
	if err != nil {
		return err
	}

	httpClient := http.DefaultClient
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	client, err := sirius.NewClient(httpClient, siriusURL)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           server.New(logger, client, tmpls, prefix, siriusPublicURL, webDir),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("listen and serve error", slog.Any("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Info("signal received: ", "sig", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return server.Shutdown(tc)
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return def
}
