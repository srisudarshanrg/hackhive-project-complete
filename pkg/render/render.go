package render

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/config"
	"github.com/srisudarshanrg/hackhive-project-competition/pkg/models"
)

var a *config.AppConfig

func SetAppConfig(app *config.AppConfig) {
	a = app
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var templateCache map[string]*template.Template
	var err error

	templateCache = a.TemplateCache

	template, ok := templateCache[tmpl]
	if !ok {
		log.Println("Requested template could not be taken from template cache")
	}

	csrfToken := nosurf.Token(r)

	td.CSRFToken = csrfToken

	err = template.Execute(w, td)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	// get all files of pattern *.page.tmpl from ./template
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Println("Could not get required pages")
		return templateCache, err
	}

	for _, page := range pages {
		filename := filepath.Base(page)
		templateSet, err := template.New(filename).ParseFiles(page)
		if err != nil {
			log.Println("Could not create template from page")
			return templateCache, err
		}

		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Println("Error while getting layout files")
			return templateCache, err
		}

		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Println(err)
				return templateCache, err
			}
		}

		templateCache[filename] = templateSet
	}

	return templateCache, nil
}
