package rest

import (
	"net/http"
	

	"github.com/labstack/echo/v4"
	"project-golang/models"
	"project-golang/repository"
	
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func (h *UserHandler) GetUser(c echo.Context) error {
	email := c.QueryParam("email")
	user, err := h.Repo.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Status:  "error",
			Message: "Failed to fetch user",
		})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateIDCard(c echo.Context) error {
	req := new(models.UpdateIDCardRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	if len(req.IDCard) != 16 {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Status:  "error",
			Message: "Invalid ID KTP, must be 16 characters",
		})
	}

	err := h.Repo.UpdateIDCard(*req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Status:  "error",
			Message: "Failed to update ID card",
		})
	}

	return c.JSON(http.StatusOK, models.ResponseMessage{
		Status:  "success",
		Message: "ID card updated successfully",
	})
}
