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
	type LoginPageData struct {
		Test string
		Form types.Form
	}

	data := LoginPageData{
		Test: "This is a test",
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

	c.Render(w, r, "auth/login", data)
}

func (c *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {

	// Handle the login logic here
	// For now, just redirect to the index page
	c.Redirect(w, r, "/", http.StatusSeeOther)
}
