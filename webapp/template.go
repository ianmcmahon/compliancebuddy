package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
)

var templates map[string]*template.Template

var baseTemplates []string = []string{"_head.tmpl", "_nav.tmpl", "_scripts.tmpl", "_base.tmpl"}

var ErrTemplateDoesNotExist = errors.New("The template does not exist.")

func rehome(files []string) []string {
	for i, v := range files {
		files[i] = fmt.Sprintf("templates/%s", v)
	}
	return files
}

func init() {
	templates = map[string]*template.Template{}

	for k, files := range templateMap {
		if files[len(files)-1] != "_base.tmpl" {
			files = append(files, baseTemplates...)
		}
		templates[k] = template.Must(template.ParseFiles(rehome(files)...))
	}
}

func renderTemplate(w io.Writer, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return ErrTemplateDoesNotExist
	}

	return tmpl.ExecuteTemplate(w, "_base.tmpl", data)
}
