package lib

import (
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"
	"time"
)

type StartResponse struct {
	// Can be 'already_running', 'in_startup', 'starting', 'creating' or 'failure'
	Status string `json:"status"`
}

type StatusResponse struct {
	// Can be 'active', 'startup', 'startup_error' or 'off'
	Status string `json:"status"`
	// Must be smaller or equal to ProgressMax
	Progress     int       `json:"progress"`
	ProgressMax  int       `json:"progress_max"`
	Name         string    `json:"name"`
	Ip           string    `json:"ip"`
	OnlinePlayer int       `json:"online_player"`
	LastOnline   time.Time `json:"last_online"`
}

func Start(port int) {
	webHandler := NewWebHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/start/", webHandler.handleStart)
	mux.HandleFunc("/status/", webHandler.handleStatus)

	log.Printf("Starting web server on port %v\n", port)
	handler := cors.AllowAll().Handler(mux)
	err := http.ListenAndServe(":"+strconv.Itoa(port), handler)

	if err != nil {
		log.Fatalf("couldn't start web server on port %v: %v\n", port, err)
	}
}
