package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/config"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/driver"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/handlers"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/render"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Could not create template cache.")
	}

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	app.TemplateCache = templateCache
	app.Session = session

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
