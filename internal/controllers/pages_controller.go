package controllers

import (
	"gohst/internal/render"
	"log"
	"net/http"

	"gohst/internal/session"
)

type PagesController struct {
	*BaseController
}

func NewPagesController() *PagesController {
    pages := &PagesController{
        BaseController: NewBaseController(),
    }

    return pages
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	sm := session.NewSessionManager()
	_, sessionId := sm.GetSession(r)
	isAuthorized, _ := sm.GetValue(sessionId, "Authorized")
	log.Println("Is authorized:", isAuthorized)
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
