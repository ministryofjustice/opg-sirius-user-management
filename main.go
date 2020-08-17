package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/ministryofjustice/opg-sirius-user-management/internal/server"
	"github.com/ministryofjustice/opg-sirius-user-management/internal/sirius"
)

func main() {
	var port, webDir, siriusURL string

	flag.StringVar(&port, "port", "9000", "Port to run on")
	flag.StringVar(&webDir, "web", "web", "Path to the 'web' directory")
	flag.StringVar(&siriusURL, "sirius", "http://localhost:9001", "URL for Sirius")
	flag.Parse()

	templates, err := template.ParseGlob(webDir + "/template/*.gotmpl")
	if err != nil {
		log.Fatal(err)
	}

	client, err := sirius.NewClient(http.DefaultClient, siriusURL)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(webDir, client, templates)

	log.Println("Running at :" + port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Println(err)
	}
}
