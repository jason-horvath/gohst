package controllers

import (
	"gohst/internal/render"
	"gohst/internal/session"
	"net/http"

	"github.com/a-h/templ"
)

type BaseController struct {
	View *render.View
}

func NewBaseController() *BaseController {
	return &BaseController{
		View: render.NewView(),
	}
}

func (c *BaseController) Render(w http.ResponseWriter, r *http.Request, page render.Page) {
	c.View.Render(w, r, page)
}

func (c *BaseController) RenderPartial(w http.ResponseWriter, r *http.Request, component templ.Component) {
	c.View.RenderPartial(w, r, component)
}

func (c *BaseController) Redirect(w http.ResponseWriter, r *http.Request, urlStr string, statusCode int) {
	http.Redirect(w, r, urlStr, statusCode)
}

func (c *BaseController) SetError(r *http.Request, message string) {
	sess := session.FromContext(r.Context())
	sess.SetFlash("error", message)
}

// SetTitle sets the page title for the current request
func (c *BaseController) SetTitle(title string) {
	c.View.SetTitle(title)
}

// SetMeta sets the page metadata for the current request.
func (c *BaseController) SetMeta(meta *render.PageMeta) {
	c.View.SetMeta(meta)
}

func (c *BaseController) JSON(w http.ResponseWriter, status int, data interface{}) {
	render.JSON(w, status, data)
}
