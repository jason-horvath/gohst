package controllers

type AuthController struct {
	*BaseController
}

func NewAuthController() *AuthController {
    auth := &AuthController{
        BaseController: NewBaseController(),
    }

    return auth
}
