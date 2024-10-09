package main

import (
	"encoding/json"
	"net/http"
)

type CommandPayload struct {
	Command string `json:"command"`
}

func (a *serverApplication) cpuCommandHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload CommandPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Command == "" {
		http.Error(w, "Command is required", http.StatusBadRequest)
		return
	}

	a.app.WorkerQueue <- payload.Command
	a.logger.Printf("Command enqueued: %s", payload.Command)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Command enqueued successfully"))
}
