package template

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

type MissingError string

func (e MissingError) Error() string {
	return fmt.Sprintf("no template exists with name '%s'", string(e))
}

type Templates map[string]*template.Template

type Template func(io.Writer, interface{}) error

// Get name template from the Templates collection.
func (t Templates) Get(name string) Template {
	tmpl := t[name]

	if tmpl == nil {
		return func(wr io.Writer, data interface{}) error {
			return MissingError(name)
		}
	}

	return func(wr io.Writer, data interface{}) error {
		return tmpl.ExecuteTemplate(wr, "page", data)
	}
}

// Parse each html/template in templateGlob with the templates in layoutGlob,
// providing the custom funcs.
func Parse(layoutGlob, templateGlob string, funcs template.FuncMap) (Templates, error) {
	layouts, err := template.New("").Funcs(funcs).ParseGlob(layoutGlob)
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(templateGlob)
	if err != nil {
		return nil, err
	}

	tmpls := map[string]*template.Template{}
	for _, file := range files {
		clone, err := layouts.Clone()
		if err != nil {
			return nil, err
		}

		tmpl, err := clone.ParseFiles(file)
		if err != nil {
			return nil, err
		}

		tmpls[filepath.Base(file)] = tmpl
	}

	return tmpls, nil
}
