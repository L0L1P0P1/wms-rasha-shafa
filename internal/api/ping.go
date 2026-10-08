package api

import (
	"encoding/json"
	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Pong string `json:"pong"`
	}{Pong: "pong"})
}
