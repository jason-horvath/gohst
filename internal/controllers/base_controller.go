package controllers

import (
	"gohst/internal/render"
	"html/template"
	"net/http"
)

type BaseController struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Templates   *template.Template
	html		*render.Html
}

func (c *BaseController) Init(w http.ResponseWriter, r *http.Request) {
	c.Writer = w
	c.Request = r
}

func (c *BaseController) Redirect(w http.ResponseWriter, urlStr string, statusCode int) {
    http.Redirect(w, c.Request, urlStr, statusCode)
}
