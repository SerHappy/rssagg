package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/db"
	"github.com/serhappy/rssagg/internal/middleware"
	"github.com/serhappy/rssagg/internal/response"
)

type FeedFollowHandler struct {
	App *app.App
}

func NewFeedFollowHandler(app *app.App) *FeedFollowHandler {
	return &FeedFollowHandler{
		App: app,
	}
}

func (h *FeedFollowHandler) CreateFeedFollow(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("Couldn't decode parameters: %v", err))
		return
	}
	feedFollow, err := h.App.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    params.FeedID,
		UserID:    user.ID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create feed follow %v", err))
		return
	}
	response.JSON(w, http.StatusCreated, feedFollow)
}

func (h *FeedFollowHandler) GetUserFeedFollows(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	feedFollows, err := h.App.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get user feed follows %v", err))
		return
	}

	response.JSON(w, http.StatusOK, feedFollows)
}

func (h *FeedFollowHandler) DeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("Couldn't parse feed follow id %v", err))
		return
	}

	err = h.App.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete feed follow %v", err))
		return
	}

	response.JSON(w, http.StatusOK, struct{}{})
}
