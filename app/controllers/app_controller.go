package controllers

import (
	"gohst/internal/controllers"
)

type AppController struct {
	*controllers.BaseController
}

func NewAppController() *AppController {
	app := &AppController{
		BaseController: controllers.NewBaseController(),
	}

	return app
}
