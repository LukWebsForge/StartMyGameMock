package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type WebHandler struct {
	status      string
	progress    int
	progressMax int
}

func NewWebHandler() *WebHandler {
	return &WebHandler{
		status: "off",
	}
}

func (web *WebHandler) handleStart(writer http.ResponseWriter, request *http.Request) {
	response := ""

	switch web.status {
	case "active":
		response = "already_running"
		break
	case "startup":
		response = "in_startup"
		break
	}

	if response != "" {
		jsonResponse(writer, StartResponse{Status: response})
		return
	}

	web.progress = 0
	web.progressMax = 6
	web.status = "startup"

	log.Println("Starting the virtual game server")

	// code to test the startup_progress
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			web.progress++
			log.Printf("Virtual start progress: %v / %v\n", web.progress, web.progressMax)
			if web.progress >= web.progressMax {
				ticker.Stop()

				log.Println("Virtual server started")
				web.status = "active"
				time.Sleep(time.Minute)
				if web.status == "active" {
					web.status = "off"
					log.Println("Stopped the virtual game server")
				}
			}
		}
	}()

	// code to test the startup_error
	/* go func() {
		time.Sleep(5 * time.Second)
		log.Println("Virtual start error")
		web.progress = 3
		web.status = "startup_error"
	}() */

	jsonResponse(writer, StartResponse{Status: "creating"})
}

func (web *WebHandler) handleStatus(writer http.ResponseWriter, request *http.Request) {
	status := StatusResponse{
		Status:      web.status,
		Progress:    web.progress,
		ProgressMax: web.progressMax,
	}

	if web.status == "active" {
		status.Name = "TTT oder nix"
		status.Ip = "177.35.26.86"
		status.OnlinePlayer = rand.Intn(6)
		status.LastOnline = time.Now().Add(-time.Duration(rand.Intn(10)) * time.Minute)
	} else {
		status.Name = ""
		status.Ip = ""
		status.OnlinePlayer = 0
		status.LastOnline = time.Time{}
	}

	jsonResponse(writer, status)
}

func jsonResponse(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(writer, err.Error())
	}
}
