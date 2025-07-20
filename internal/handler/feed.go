package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/db"
	"github.com/serhappy/rssagg/internal/middleware"
	"github.com/serhappy/rssagg/internal/response"
)

type FeedHandler struct {
	App *app.App
}

func NewFeedHandler(app *app.App) *FeedHandler {
	return &FeedHandler{
		App: app,
	}
}

func (h *FeedHandler) CreateFeed(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("Couldn't decode parameters: %v", err))
		return
	}
	feed, err := h.App.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create feed %v", err))
		return
	}
	response.JSON(w, http.StatusCreated, feed)
}

func (h *FeedHandler) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.App.DB.GetFeeds(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Couldn't get feeds")
		return
	}

	response.JSON(w, http.StatusOK, feeds)
}
