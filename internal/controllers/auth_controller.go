package controllers

import (
	"net/http"

	"gohst/internal/auth"
	"gohst/internal/session"
	"gohst/internal/types"
	"gohst/internal/utils"
	"gohst/internal/validation"
)
type AuthController struct {
	*BaseController
}

func NewAuthController() *AuthController {
    auth := &AuthController{
        BaseController: NewBaseController(),
    }
	auth.view.SetLayout("layouts/auth")
    return auth
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	sess := session.FromContext(r.Context())
	type LoginPageData struct {
		Test string
		Form types.Form
	}

	emailValue, _ := sess.PeekOld("email")

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
						Value: utils.StringOr(emailValue, ""),
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

// Handle the login info that is submitted from the form
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
        sess.SetFlash("login_error", "Email and password are required")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }


	if !validation.IsEmail(email) {
		sess.SetFlash("login_error", "Invalid email format")
		c.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

    // Find user in database
	user, err := auth.Login(sess, email, password)
    if err != nil {
        sess.SetFlash("login_error", "Invalid email or password")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Verify password
	passwordOk, _ := utils.CheckPassword(password, user.PasswordHash)
    if (!passwordOk) {
        sess.SetFlash("login_error", "Invalid email or password")
        c.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Redirect to dashboard or home page
    c.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	sess := session.FromContext(r.Context())
	type LoginPageData struct {
		Test string
		Form types.Form
	}

	firstName, _ := sess.PeekOld("first_name")
	lastName, _ := sess.PeekOld("last_name")
	emailValue, _ := sess.PeekOld("email")
	emailConfrimValue, _ := sess.PeekOld("email_confirm")

	data := LoginPageData{
		Test: "This is a test",
		Form: types.Form{
			Method: "POST",
			Action: "/login", // Adjust as needed
			Fieldset: types.Fieldset{
				"first_name": types.Field{
                Input: types.Text{
                    Name:        "first_name",
                    Type:        "text",
                    ID:          "first_name",
                    Placeholder: "Your first name",
                    Value:       utils.StringOr(firstName, ""),
                },
                Label: types.Label{For: "first_name", Text: "First Name"},
				},
				"last_name": types.Field{
					Input: types.Text{
						Name:        "last_name",
						Type:        "text",
						ID:          "last_name",
						Placeholder: "Your last name",
						Value:       utils.StringOr(lastName, ""),
					},
					Label: types.Label{For: "last_name", Text: "Last Name"},
				},
				"email": types.Field{
					Input: types.Text{
						Name: "email",
						Type: "email",
						ID: "email",
						Placeholder: "Enter your email.",
						Value: utils.StringOr(emailValue, ""),
					},
					Label: types.Label{For: "email", Text: "Email"},
				},
				"email_confirm": types.Field{
					Input: types.Text{
						Name: "email_confirm",
						Type: "email_confirm",
						ID: "email_confirm",
						Placeholder: "Confirm your email.",
						Value: utils.StringOr(emailConfrimValue, ""),
					},
					Label: types.Label{For: "email", Text: "Confirm Email"},
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
				"password_confirm": types.Field{
					Input: types.Text{
						Name: "password_confirm",
						Type: "password_confirm",
						ID: "password_confirm",
						Placeholder: "Confirm your password.",
					},
					Label: types.Label{For: "password", Text: "Confirm Password"},
				},
			},
			Buttons: map[string]types.Button{
				"submit": {
					Type: "submit",
					Text: "Register",
				},
			},
		},
	}

	c.Render(w, r, "auth/register", data)
}

func (c *AuthController) HandleRegister(w http.ResponseWriter, r *http.Request) {
    sess := session.FromContext(r.Context())

    // Parse the form data
    if err := r.ParseForm(); err != nil {
        c.SetError(r, "Failed to parse form data")
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Get form values
    email := r.FormValue("email")
    emailConfirm := r.FormValue("email_confirm")
    firstName := r.FormValue("first_name") // Add these fields to your form
    lastName := r.FormValue("last_name")   // Add these fields to your form
    password := r.FormValue("password")
    passwordConfirm := r.FormValue("password_confirm")

    // Save for form repopulation
    sess.SetOld("email", email)
    sess.SetOld("email_confirm", emailConfirm)
    sess.SetOld("first_name", firstName)
    sess.SetOld("last_name", lastName)

    // Validate inputs
    if email == "" || emailConfirm == "" || password == "" || passwordConfirm == "" {
        sess.SetFlash("register_error", "All fields are required")
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Check if emails match
    if email != emailConfirm {
        sess.SetFlash("register_error", "Emails do not match")
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Check if passwords match
    if password != passwordConfirm {
        sess.SetFlash("register_error", "Passwords do not match")
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Validate email format
    if !validation.IsEmail(email) {
        sess.SetFlash("register_error", "Invalid email format")
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Validate password strength
    if len(password) < 8 {
        sess.SetFlash("register_error", "Password must be at least 8 characters")
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Register the user
    err := auth.Register(email, firstName, lastName, password)
    if err != nil {
        sess.SetFlash("register_error", err.Error())
        c.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    // Set success message
    sess.SetFlash("login_success", "Registration successful! You can now log in.")

    // Redirect to login page
    c.Redirect(w, r, "/login", http.StatusSeeOther)
}
