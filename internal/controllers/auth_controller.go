package controllers

import (
	"net/http"
)
type AuthController struct {
	*BaseController
}

func NewAuthController() *AuthController {
    auth := &AuthController{
        BaseController: NewBaseController(),
    }

    return auth
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	c.Init(w, r)
	c.html.RenderView(w, "auth/login.html")
}
