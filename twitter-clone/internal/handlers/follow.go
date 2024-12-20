package handlers

import (
	"net/http"
	"strconv"
	"twitter-clone/internal/services"

	"github.com/labstack/echo/v4"
)

type FollowHandler struct {
	followService *services.FollowService
}

func NewFollowHandler(followService *services.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

func (h *FollowHandler) FollowUser(c echo.Context) error {
	followerID := c.Get("user_id").(uint)
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if err := h.followService.FollowUser(followerID, uint(followingID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *FollowHandler) UnfollowUser(c echo.Context) error {
	followerID := c.Get("user_id").(uint)
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	if err := h.followService.UnfollowUser(followerID, uint(followingID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *FollowHandler) GetFollowers(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	followers, err := h.followService.GetFollowers(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, followers)
}

func (h *FollowHandler) GetFollowing(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	following, err := h.followService.GetFollowing(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, following)
}

func (h *FollowHandler) CheckFollowStatus(c echo.Context) error {
	followerID := c.Get("user_id").(uint)
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	isFollowing := h.followService.IsFollowing(followerID, uint(followingID))
	return c.JSON(http.StatusOK, map[string]bool{"is_following": isFollowing})
}
