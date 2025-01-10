package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/HsiaoCz/go-master/bunt/db"
	"github.com/HsiaoCz/go-master/bunt/pkgs"
	"github.com/HsiaoCz/go-master/bunt/types"
	"github.com/google/uuid"
)

func HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create a new post
	// ...
	var create_post types.CreatePost
	if err := json.NewDecoder(r.Body).Decode(&create_post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post := &types.Posts{
		PostID:    uuid.New().String(),
		UserID:    create_post.UserID,
		Content:   create_post.Content,
		ImageUrl:  create_post.ImageUrl,
		VideoUrl:  create_post.VideoUrl,
		Title:     create_post.Title,
		Caption:   create_post.Caption,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	session, ok := r.Context().Value(types.CtxSessionKey).(*types.Sessions)
	if !ok {
		http.Error(w, "please login", http.StatusForbidden)
		return
	}

	location, err := pkgs.GetLocationByIP(session.IpAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	post.Location = location

	_, err = db.Get().NewInsert().Model(post).Exec(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"message": "post created successfully", "status": http.StatusCreated})
}
