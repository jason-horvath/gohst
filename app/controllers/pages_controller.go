package controllers

import (
	"net/http"

	appConfig "gohst/app/config"
	"gohst/internal/session"
)

type PagesController struct {
	*AppController
}

func NewPagesController() *PagesController {
    pages := &PagesController{
        AppController: NewAppController(),
    }

    return pages
}

func (c *PagesController) Index(w http.ResponseWriter, r *http.Request) {
	sess := session.FromContext(r.Context())
	username, _ := sess.Get("Username")

	data := map[string]interface{}{
		"SessionID":    sess.ID(),
		"Username":     username,
		"AppName":      appConfig.App.Name,
		"AppVersion":   appConfig.App.Version,
		"IsProduction": appConfig.App.IsProduction(),
		"Features":     appConfig.App.Features,
	}

	c.Render(w, r, "pages/index", data)
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
    c.Render(w, r, "pages/404")
}
