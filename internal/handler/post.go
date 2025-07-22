package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/db"
	"github.com/serhappy/rssagg/internal/middleware"
	"github.com/serhappy/rssagg/internal/response"
)

const (
	defaultLimit = 10
)

type PostHandler struct {
	App *app.App
}

func NewPostHandler(app *app.App) *PostHandler {
	return &PostHandler{
		App: app,
	}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	limit := defaultLimit
	queryLimitStr := r.URL.Query().Get("limit")
	if queryLimit, err := strconv.Atoi(queryLimitStr); err == nil {
		limit = queryLimit
	}

	posts, err := h.App.DB.GetNewestPostsForUser(r.Context(), db.GetNewestPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get posts for user %v", err))
		return
	}

	response.JSON(w, http.StatusOK, databasePostsToPosts(posts))
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(post db.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: nullStringToStringPtr(post.Description),
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}

func databasePostsToPosts(posts []db.Post) []Post {
	result := make([]Post, len(posts))
	for i, post := range posts {
		result[i] = databasePostToPost(post)
	}
	return result
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}
