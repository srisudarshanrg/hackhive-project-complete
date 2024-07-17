package main

import (
	"log"
	"net/http"

	"github.com/srisudarshanrg/hackhive-project-competition/pkg/config"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/driver"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/handlers"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/render"
)

const portNumber = ":8000"

var app config.AppConfig

func main() {
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Could not create template cache.")
	}

	app.TemplateCache = templateCache

	repo := handlers.SetAppConfigHandler(&app)
	handlers.NewHandlers(repo)
	render.SetAppConfig(&app)

	// create database
	db, err := driver.CreateDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	handlers.DatabaseAccess(db)
	defer db.Close()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal("Could not serve application", err)

}
