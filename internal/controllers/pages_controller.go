package controllers

import (
	"gohst/internal/render"
	"log"
	"net/http"

	"gohst/internal/models"
	"gohst/internal/session"
	"gohst/internal/utils"
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
	hashed, _ := utils.HashPassword("Test1234!") // Example of hashing a password
	userModel := models.NewUserModel()
	user, err := userModel.FirstOf("SELECT * FROM users WHERE email = $1", "admin@example.com")

	if err != nil {
		log.Println("Error fetching user:", err)
	}
	log.Println("User:", user)
	data := map[string]interface{}{
		"SessionID": sess.ID(),
		"Username":  username,
		"Hashed":  "This is a password: " + hashed,
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

    render.JSON(w, response)
}

// NotFound handles 404 errors
func (c *PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    c.Render(w, r, "pages/404")
}
