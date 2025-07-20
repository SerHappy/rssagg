package handler

import (
	"net/http"

	"github.com/serhappy/rssagg/internal/response"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
