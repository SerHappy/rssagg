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

type UserHandler struct {
	App *app.App
}

func NewUserHandler(app *app.App) *UserHandler {
	return &UserHandler{
		App: app,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, fmt.Sprintf("Couldn't decode parameters: %v", err))
		return
	}
	user, err := h.App.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	response.JSON(w, http.StatusCreated, user)

}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUser(r.Context())
	if !ok {
		response.Error(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	response.JSON(w, http.StatusOK, user)
}
