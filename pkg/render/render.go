package render

import (
	"bytes"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Renders web files

func RenderTemplate(w http.ResponseWriter, gohtml string) {

	parserdTemplate, _ := template.ParseFiles("templates/" + gohtml)
	err := parserdTemplate.Execute(w, nil)
	if err != nil {
		log.Println("Error parsing template", err)
		return
	}

}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

// RenderTemplateN similar to the previous but renders go template
func RenderTemplateN(w http.ResponseWriter, gohtml string, td *models.TemplateData) {

	var templateCache map[string]*template.Template
	if app.UseCache {
		// create a template cache for pages
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
		/*	if err != nil {
			log.Println("Error getting template cache", err)
			log.Fatal(err) //this exits the application
			//return
		}*/
	}

	// create a template cache for pages
	/*templateCache, err := CreateTemplateCache()
	if err != nil {
		log.Println("Error getting template cache", err)
		log.Fatal(err) //this exits the application
		//return
	}*/
	// get specific gohtml template
	t, ok := templateCache[gohtml] //ok for if it does not exist
	if !ok {
		log.Fatal("Specific path not found")
	}
	// since the template won't be read from disk we use a buffer
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, td)
	// writes to the response html
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error in writing template to browser", err)
	}

}

// map of functions to do some processing and pass to the template
var functions = template.FuncMap{}

// CreateTemplateCache creates a map of templates
func CreateTemplateCache() (map[string]*template.Template, error) {

	// this stores the template in a cache
	myCache := map[string]*template.Template{}

	//get all template with .page.gohtml
	pages, err := filepath.Glob("templates/*.page.gohtml")
	if err != nil {
		return myCache, err
	}
	//loop through all the pages
	for _, page := range pages {
		// get the pages base/name
		name := filepath.Base(page)

		// pass the template file
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// check if the template matches any layout
		matches, err := filepath.Glob("templates/*.layout.gohtml")
		if err != nil {
			return myCache, err
		}

		// check if any layout is found and pass the 'pages' into the 'layout'
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("templates/*.layout.gohtml")
			if err != nil {
				return myCache, err
			}
		}
		// store the cache templates in a map
		myCache[name] = ts

	}

	return myCache, nil
}
