package handler

import (
	"net/http"

	"github.com/serhappy/rssagg/internal/response"
)

func Error(w http.ResponseWriter, r *http.Request) {
	response.Error(w, http.StatusInternalServerError, "Internal Server Error")
}
