package controllers

import (
	"gohst/internal/render"
	"gohst/internal/session"
	"gohst/internal/utils"
	"html/template"
	"net/http"
)

type BaseController struct {
	Templates   *template.Template
	View        *render.View
}

func NewBaseController() *BaseController {

	view := render.NewView()
	base := &BaseController{
		View: view,
	}

	return base
}

func (c *BaseController) Render(w http.ResponseWriter, r *http.Request, viewName string, data ...interface{}) {
	useData := utils.StructSafe(data)
    c.View.Render(w, r, viewName, useData)
}

func (c *BaseController) Redirect(w http.ResponseWriter, r *http.Request, urlStr string, statusCode int) {
    http.Redirect(w, r, urlStr, statusCode)
}

func (c *BaseController) SetError(r *http.Request, message string) {
    sess := session.FromContext(r.Context())
    sess.SetFlash("error", message)
}
