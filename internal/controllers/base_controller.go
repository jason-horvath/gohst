package controllers

import (
	"gohst/internal/render"
	"gohst/internal/utils"
	"html/template"
	"net/http"
)

type BaseController struct {
	Templates   *template.Template
	view		*render.View
}

func NewBaseController() *BaseController {

	view := render.NewView()
	base := &BaseController{
		view: view,
	}

	return base
}

func (c *BaseController) Render(w http.ResponseWriter, r *http.Request, viewName string, data ...interface{}) {
	useData := utils.StructEmpty(data)
    c.view.Render(w, r, viewName, useData)
}


func (c *BaseController) Redirect(w http.ResponseWriter, r *http.Request, urlStr string, statusCode int) {
    http.Redirect(w, r, urlStr, statusCode)
}
