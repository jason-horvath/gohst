package controllers

import (
	"net/http"

	"gohst/views/pages"
)

type PagesController struct {
	*AppController
}

func NewPagesController() *PagesController {
    p := &PagesController{
        AppController: NewAppController(),
    }

    return p
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	c.Render(w, r, pages.IndexPage())
}

func (c *PagesController) Post(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

    response := struct {
        ID      string `json:"id"`
        Message string `json:"message"`
    }{
        ID:      id,
        Message: "This is post " + id,
    }

    c.JSON(w, http.StatusOK, response)
}

// NotFound handles 404 errors
func (c *PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    c.Render(w, r, pages.NotFoundPage())
}
