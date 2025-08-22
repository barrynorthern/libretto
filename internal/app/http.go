package app

import (
	"encoding/json"
	"net/http"
)

// RegisterHTTP mounts JSON endpoints used by the UI.
type StoreReader interface {
	ListScenes(r *http.Request) any
}

func RegisterHTTP(mux *http.ServeMux, gw StoreReader) {
	mux.HandleFunc("/api/scenes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		list := gw.ListScenes(r)
		_ = json.NewEncoder(w).Encode(list)
	})
}
