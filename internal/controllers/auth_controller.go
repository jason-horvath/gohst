package controllers

import (
	"net/http"

	"gohst/internal/models"
	"gohst/internal/session"
	"gohst/internal/types"
	"gohst/internal/utils"
)
type AuthController struct {
	*BaseController
}

func NewAuthController() *AuthController {
    auth := &AuthController{
        BaseController: NewBaseController(),
    }
	auth.view.SetLayout("layout/auth")
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
						ID: "email",
						Placeholder: "Enter your email.",
					},
					Label: types.Label{For: "email", Text: "Email"},
				},
				"password": types.Field{
					Input: types.Text{
						Name: "password",
						Type: "password",
						ID: "password",
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
	sess := session.FromContext(r.Context())

	// Parse the form data
    if err := r.ParseForm(); err != nil {
        c.SetError(r, "Failed to parse form data")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Get email and password from the form
    email := r.FormValue("email")
    password := r.FormValue("password")
	sess.SetOld("email", email)

    // Validate input
    if email == "" || password == "" {
        c.SetError(r, "Email and password are required")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Find user in database
    userModel := models.NewUserModel()
    user, err := userModel.FindByEmail(email)

    if err != nil {
        c.SetError(r, "Invalid email or password")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Verify password
	passwordOk, _ := utils.CheckPassword(password, user.PasswordHash)
    if (!passwordOk) {
        c.SetError(r, "Invalid email or password")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    sess.Set("UserID", user.ID)
    sess.Set("Email", user.Email)

    // Redirect to dashboard or home page
    c.Redirect(w, r, "/", http.StatusSeeOther)
}
