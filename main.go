package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/server"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/env"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/logging"
	"github.com/ministryofjustice/opg-sirius-user-management/tbd/template"
)

func main() {
	logger := logging.New(os.Stdout, "opg-sirius-user-management")

	port := env.Get("PORT", "8080")
	webDir := env.Get("WEB_DIR", "web")
	siriusURL := env.Get("SIRIUS_URL", "http://localhost:9001")
	siriusPublicURL := env.Get("SIRIUS_PUBLIC_URL", "")
	prefix := env.Get("PREFIX", "")

	tmpls, err := template.Parse(webDir+"/template/layout/*.gotmpl", webDir+"/template/*.gotmpl", map[string]interface{}{
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
	})
	if err != nil {
		logger.Fatal(err)
	}

	client, err := sirius.NewClient(http.DefaultClient, siriusURL)
	if err != nil {
		logger.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: server.New(logger, client, tmpls, prefix, siriusPublicURL, webDir),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Print("Running at :" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Print("signal received: ", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(tc); err != nil {
		logger.Print(err)
	}
}
