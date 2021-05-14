package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

var PathToTemplates = "./../../templates"

// // RenderTemplate using http template
func Template(w http.ResponseWriter, tmpl string) error {

	var tc map[string]*template.Template

	tc, err := CreateTemplateCache()
	if err != nil {
		log.Println(err)
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Println("Cant get template from cache")
		return errors.New("cant get template from cache")

	}

	buff := new(bytes.Buffer)

	err = t.Execute(buff, nil)
	if err != nil {
		log.Println(err)
		log.Fatal("execution error")
	}


	_, err = buff.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}


	return nil
}

// Creates a template cache as a map

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}


	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", PathToTemplates))
	if err != nil {
		return myCache, err
	}

	fmt.Println(pages)

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}


		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", PathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", PathToTemplates))
			if err != nil {
				fmt.Println("tojto6", err)
			}

		}

		myCache[name] = ts

	}

	return myCache, err
}
