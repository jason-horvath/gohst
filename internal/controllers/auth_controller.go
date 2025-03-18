package controllers

import (
	"net/http"

	"gohst/internal/types"
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

	type LoginPageData struct {
		Test string
		Label types.Label
		Form types.Form
	}

	data := LoginPageData{
		Test: "This is a test",
		Label: types.Label{
			For: "email",
			Text: "Email Address",
		},
		Form: types.Form{
			Method: "POST",
			Action: "/login", // Adjust as needed
			Fieldset: types.Fieldset{
				"email": types.Field{
					Input: types.Text{
						Name: "email",
						Type: "email",
						Placeholder: "Enter your email.",
					},
					Label: types.Label{For: "email", Text: "Email"},
				},
				"password": types.Field{
					Input: types.Text{
						Name: "password",
						Type: "password",
						Placeholder: "Enter your password.",
					},
					Label: types.Label{For: "password", Text: "Password"},
				},
			},
			Buttons: map[string]types.Button{
				"submit": {
					Type: "submit",
					Text: "Login",
				},
			},
		},
	}

	c.view.Render(w, "auth/login", data)
}
