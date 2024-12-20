package handlers

import (
	"net/http"
	"twitter-clone/internal/services"

	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) UpdateBackground(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	// Get the file from the request
	file, err := c.FormFile("background")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no file uploaded")
	}

	// Check file size (e.g., 5MB limit)
	if file.Size > 5*1024*1024 {
		return echo.NewHTTPError(http.StatusBadRequest, "file too large")
	}

	backgroundURL, err := h.profileService.UpdateBackground(userID, file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"background_url": backgroundURL,
	})
}

type UpdateProfileRequest struct {
	Bio string `json:"bio"`
}

func (h *ProfileHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.profileService.UpdateProfile(userID, req.Bio); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
