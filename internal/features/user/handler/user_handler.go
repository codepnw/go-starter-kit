package userhandler

import (
	"net/http"

	"github.com/codepnw/go-starter-kit/internal/errs"
	"github.com/codepnw/go-starter-kit/internal/features/user"
	userservice "github.com/codepnw/go-starter-kit/internal/features/user/service"
	"github.com/codepnw/go-starter-kit/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service userservice.UserService
}

func NewUserHandler(service userservice.UserService) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) Register(c *gin.Context) {
	req := new(RegisterReq)

	if err := c.ShouldBindJSON(req); err != nil {
		response.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	input := &user.User{
		Email:    req.Email,
		Password: req.Password,
	}
	resp, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		switch err {
		case errs.ErrEmailAlreadyExists:
			response.ResponseError(c, http.StatusBadRequest, err)
		default:
			response.ResponseError(c, http.StatusInternalServerError, err)
		}
		return
	}

	response.ResponseSuccess(c, http.StatusCreated, resp)
}

func (h *userHandler) Login(c *gin.Context) {
	req := new(LoginReq)

	if err := c.ShouldBindJSON(req); err != nil {
		response.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case errs.ErrInvalidEmailOrPassword:
			response.ResponseError(c, http.StatusBadRequest, err)
		default:
			response.ResponseError(c, http.StatusInternalServerError, err)
		}
		return
	}
	
	response.ResponseSuccess(c, http.StatusOK, resp)
}
