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

var PathToTemplates = "./../../templates" // za testing, ker nismo v rootu

// // RenderTemplate using http template
func Template(w http.ResponseWriter, tmpl string) error {

	var tc map[string]*template.Template

	tc, err := CreateTemplateCache() // prevermo ce je ze template cache, ce ni ga ustvari
	if err != nil {
		log.Println(err)
	}

	t, ok := tc[tmpl] // map checking ce sploh obstaja "ok", pogledamo za vsak about.page in home.page pa vstavmo v map
	if !ok {
		// log.Fatal("could not get template from template cache")
		log.Println("Cant get template from cache")
		return errors.New("cant get template from cache") // lahko zbrisemo, ker v testu zjebe ko testiramo za nonexistent template

	}

	buff := new(bytes.Buffer) // da sprejme info iz *template.Template

	err = t.Execute(buff, nil)
	if err != nil {
		log.Println(err)
		log.Fatal("execution error")
	}

	// execute the template, store the value in buffer

	_, err = buff.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl) // ->> tega ne rabmo vec, pisemo gor direkt v buf
	// err := parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("Error parsing template", err)
	// 	return
	// }
	return nil
}

// Creates a template cache as a map

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{} // uporabmo template package, tuki ze {} za template, ker definiramo

	// cache uporabmo, ker ga bomo gor uporabljal za vstavljat template home pa about

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", PathToTemplates)) // izbere vse kar je .page.html v templates
	if err != nil {
		return myCache, err
	}

	fmt.Println(pages)

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page) // damo functions ker bomo slej ko prej uporabljali funkcije
		if err != nil {
			return myCache, err
		}

		// zgornji ts dejansko ustvari prazni template za npr. home.page, spodaj pa prazni ts uproabmo da damo layout.html notri,
		// in potem je ta template poln informacij o page name, layout template...

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", PathToTemplates)) // zdej pa iscemo layout file
		if err != nil {
			fmt.Println("tojto5", err)
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", PathToTemplates))
			if err != nil {
				fmt.Println("tojto6", err)
				return myCache, err
			}

		}

		myCache[name] = ts

	}

	return myCache, err
}
