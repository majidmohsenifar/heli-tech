package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/majidmohsenifar/heli-tech/gateway-service/service/user"
)

type UserHandler struct {
	userService *user.Service
	validate    *validator.Validate
}

// This endpoint allows user to register
//	@Summary		register user
//	@Description	allows user to register
//	@Tags			User
//	@ID				register
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			params	body		user.RegisterParams	false	"Register-Params"
//	@Success		200		{object}	ResponseSuccess
//	@Failure		400		{object}	ResponseFailure
//	@Failure		500		{object}	ResponseFailure
//	@Router			/api/v1/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	params := user.RegisterParams{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Bad Request: "+err.Error(),
		)
		return
	}
	err = h.validate.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Invalid Request: "+err.Error())
		return
	}
	err = h.userService.Register(c, params)
	if err != nil {
		MakeErrorResponseWithoutCode(c.Writer, err)
		return
	}
	MakeSuccessResponse(c.Writer, nil, "successfully registered")
}

// This endpoint allows user to login
//	@Summary		login user
//	@Description	allows user to login
//	@Tags			User
//	@ID				login
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			params	body		user.RegisterParams	false	"Register-Params"
//	@Success		200		{object}	ResponseSuccess
//	@Failure		400		{object}	ResponseFailure
//	@Failure		500		{object}	ResponseFailure
//	@Router			/api/v1/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	params := user.LoginParams{}
	err := c.ShouldBindJSON(&params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Bad Request: "+err.Error(),
		)
		return
	}
	err = h.validate.Struct(params)
	if err != nil {
		MakeErrorResponseWithCode(
			c.Writer,
			http.StatusBadRequest,
			"Invalid Request: "+err.Error())
		return
	}
	res, err := h.userService.Login(c, params)
	if err != nil {
		MakeErrorResponseWithoutCode(c.Writer, err)
		return
	}
	MakeSuccessResponse(c.Writer, res, "successfully logged in")
}

func NewUserHandler(
	userService *user.Service,
	validate *validator.Validate,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validate,
	}
}
