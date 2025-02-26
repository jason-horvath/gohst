package controllers

import (
	"gohst/internal/render"
	"net/http"
)

type PagesController struct {
	*BaseController
}

func NewPagesController() *PagesController {
    html := render.NewHtml()

    pages := &PagesController{
        BaseController: &BaseController{
            html: html,
        },
    }

    return pages
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	c.Init(w, r)
	c.html.RenderView(w, "pages/index.html")
}

func (c *PagesController) About(w http.ResponseWriter, r *http.Request) {
	c.Init(w, r)
	c.html.RenderView(w, "pages/about.html")
}

func (c *PagesController) Post(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c.Init(w, r)
	render.Text(w, "This is post " + id)
}
