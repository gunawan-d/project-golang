package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"project-golang/services"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

// CreateToken akan menerima data dari request dan mengembalikan token JWT
func (h *AuthHandler) CreateToken(c echo.Context) error {
	// Bind payload dari request body
	payload := make(map[string]interface{})
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status":  "error",
			"message": "Invalid request body, required name, email, roleID",
		})
	}

	// Panggil AuthService untuk membuat token
	token, err := h.AuthService.CreateToken(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Failed to create token",
		})
	}

	// Menambahkan custom header pada response
	c.Response().Header().Set("X-Custom-Header", "CustomHeaderValue")
	c.Response().Header().Set("Authorization", "Bearer "+token)

	// Ambil nilai dari payload
	name, _ := payload["name"].(string)
	email, _ := payload["email"].(string)
	roleID, _ := payload["roleID"].(string)

	// Kembalikan token dan informasi lainnya ke client
	return c.JSON(http.StatusOK, map[string]string{
		"token":   token,
		"status":  "Success",
		"name":    name,
		"email":   email,
		"roleID":  roleID,
	})
}
