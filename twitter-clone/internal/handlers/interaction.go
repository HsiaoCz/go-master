package handlers

import (
	"net/http"
	"strconv"
	"twitter-clone/internal/services"

	"github.com/labstack/echo/v4"
)

type InteractionHandler struct {
	likeService    *services.LikeService
	retweetService *services.RetweetService
}

func NewInteractionHandler(likeService *services.LikeService, retweetService *services.RetweetService) *InteractionHandler {
	return &InteractionHandler{
		likeService:    likeService,
		retweetService: retweetService,
	}
}

func (h *InteractionHandler) LikeTweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	if err := h.likeService.LikeTweet(userID, uint(tweetID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *InteractionHandler) UnlikeTweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	if err := h.likeService.UnlikeTweet(userID, uint(tweetID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *InteractionHandler) GetTweetLikes(c echo.Context) error {
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	likes, err := h.likeService.GetTweetLikes(uint(tweetID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, likes)
}

func (h *InteractionHandler) Retweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	if err := h.retweetService.CreateRetweet(userID, uint(tweetID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusCreated)
}

func (h *InteractionHandler) UndoRetweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	if err := h.retweetService.UndoRetweet(userID, uint(tweetID)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *InteractionHandler) GetTweetRetweets(c echo.Context) error {
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	retweets, err := h.retweetService.GetTweetRetweets(uint(tweetID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, retweets)
}