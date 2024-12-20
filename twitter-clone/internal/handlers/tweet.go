package handlers

import (
	"net/http"
	"strconv"
	"twitter-clone/internal/services"

	"github.com/labstack/echo/v4"
)

type TweetHandler struct {
	tweetService *services.TweetService
}

func NewTweetHandler(tweetService *services.TweetService) *TweetHandler {
	return &TweetHandler{tweetService: tweetService}
}

type CreateTweetRequest struct {
	Content  string `json:"content" validate:"required"`
	ParentID *uint  `json:"parent_id,omitempty"`
}

func (h *TweetHandler) CreateTweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	var req CreateTweetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tweet, err := h.tweetService.CreateTweet(userID, req.Content, req.ParentID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, tweet)
}

func (h *TweetHandler) GetTweet(c echo.Context) error {
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	tweet, err := h.tweetService.GetTweet(uint(tweetID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "tweet not found")
	}

	return c.JSON(http.StatusOK, tweet)
}

func (h *TweetHandler) UpdateTweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	var req CreateTweetRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tweet, err := h.tweetService.UpdateTweet(uint(tweetID), userID, req.Content)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tweet)
}

func (h *TweetHandler) DeleteTweet(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	if err := h.tweetService.DeleteTweet(uint(tweetID), userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *TweetHandler) GetUserTweets(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	tweets, err := h.tweetService.GetUserTweets(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tweets)
}

type TimelineQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (h *TweetHandler) GetHomeTimeline(c echo.Context) error {
	userID := c.Get("user_id").(uint)

	query := new(TimelineQuery)
	if err := c.Bind(query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Set default values if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20 // Default limit
	}

	response, err := h.tweetService.GetTimelineWithMetadata(userID, query.Page, query.Limit, "home")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TweetHandler) GetUserTimeline(c echo.Context) error {
	targetUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	query := new(TimelineQuery)
	if err := c.Bind(query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Set default values if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20 // Default limit
	}

	response, err := h.tweetService.GetTimelineWithMetadata(uint(targetUserID), query.Page, query.Limit, "user")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (h *TweetHandler) GetTweetReplies(c echo.Context) error {
	tweetID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tweet ID")
	}

	query := new(TimelineQuery)
	if err := c.Bind(query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Set default values if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 20 // Default limit
	}

	replies, err := h.tweetService.GetReplies(uint(tweetID), query.Page, query.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, replies)
}
