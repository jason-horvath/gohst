package controllers

import (
	"gohst/internal/controllers"
)

// AppController provides application-specific functionality that extends the framework's BaseController.
// This is the base controller for all application controllers and should be embedded rather than used directly.
//
// Usage:
//   type AuthController struct {
//       *AppController
//   }
//
//   func NewAuthController() *AuthController {
//       return &AuthController{
//           AppController: NewAppController(),
//       }
//   }
//
// AppController adds application-level functionality on top of the framework's BaseController,
// such as custom logging, app-specific middleware, and business logic helpers.
// All application controllers should embed AppController to maintain consistency
// and gain access to both framework and application-level functionality.
type AppController struct {
	*controllers.BaseController
}

// NewAppController creates a new AppController instance.
// This should typically be called from other controller constructors rather than used directly.
func NewAppController() *AppController {
	app := &AppController{
		BaseController: controllers.NewBaseController(),
	}

	return app
}
