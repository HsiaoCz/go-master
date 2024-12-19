package handlers

import (
	"net/http"

	"github.com/blogshare/internal/models"
	"github.com/blogshare/internal/repository"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	repo *repository.ElasticsearchRepository
}

func NewPostHandler(repo *repository.ElasticsearchRepository) *PostHandler {
	return &PostHandler{repo: repo}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPost := models.NewPost(post.Title, post.Content, post.Author, post.Tags)
	if err := h.repo.CreatePost(c.Request.Context(), newPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newPost)
}

func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query is required"})
		return
	}

	posts, err := h.repo.SearchPosts(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
