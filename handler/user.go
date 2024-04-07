package handler

import (
	"net/http"

	"github.com/funukonta/crowd-funding/helper"
	"github.com/funukonta/crowd-funding/user"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *userHandler {
	return &userHandler{service: service}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	input := user.RegisterUserInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", gin.H{"errors": errors})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.service.RegisterUser(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	formatter := user.FormatUser(newUser, "token")
	response := helper.APIResponse("Account created!", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
