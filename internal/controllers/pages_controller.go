package controllers

import (
	"gohst/internal/render"
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
	sess := session.FromContext(r.Context())
	username, _ := sess.Get("Username") // Example of getting a user from the session
	data := map[string]interface{}{
		"SessionID": sess.ID(),
		"Username":  username,
	}

	c.Render(w, r, "pages/index", data)
}

func (c *PagesController) About(w http.ResponseWriter, r *http.Request) {
	c.Render(w, r, "pages/about")
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

    render.JSON(w, response)
}
