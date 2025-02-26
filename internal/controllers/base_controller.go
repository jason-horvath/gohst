package controllers

import (
	"gohst/internal/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type BaseController struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Templates   *template.Template
	viewsDir 	string
	html		*render.Html
}

func (c *BaseController) Init(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
}

func (c *BaseController) loadTemplateDir() {
	render.LoadTemplateDir(c.viewsDir)
}

func (c *BaseController) renderTemplate(w http.ResponseWriter, tmpl string) {
	// ✅ Reloads the template every time (no caching)
	tmplPath := filepath.Join("templates", tmpl)
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Println("Error loading template:", err)
		return
	}
	t.Execute(w, nil)
}
