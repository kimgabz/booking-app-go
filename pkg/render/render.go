package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/kimgabz/booking-app-go/pkg/config"
	"github.com/kimgabz/booking-app-go/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package.
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data to the templates that needs to be on every page.
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using http/template.
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	// Create variable to hold template cache.
	var tc map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the app config.
		tc = app.TemplateCache
	} else {
		// Rebuild the template cache.
		// (useful for viewing changes to template in Dev, not Production where
		// the template will not be changed often.)
		tc, _ = CreateTemplateCache()
	}

	// Get the requsted template by its string name via
	// map from the template cache.
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache.")
	}

	// Create a new buffer.
	buf := new(bytes.Buffer)

	// Add data that should be present on all templates.
	td = AddDefaultData(td)

	// Execute the template and data into the buffer.
	_ = t.Execute(buf, td)

	// Write the buffer contents to the response writer.
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}
}

// CreateTemplateCache creates a template cache as a map.
func CreateTemplateCache() (map[string]*template.Template, error) {
	// 	Create a cache to store our parsed templates.
	myCache := map[string]*template.Template{}

	// 	Find pages in template folder.
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	// 	Range through pages found in template folder.
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
