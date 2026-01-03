package routes

import (
	"encoding/json"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/version"
)

type SystemHealth struct {
	DatabaseStatus        bool   `json:"databaseStatus"`
	DatabaseStatusMessage string `json:"databaseStatusMessage"`
	Version               string `json:"version"`
}

func (routes *Routes) Health(w http.ResponseWriter, r *http.Request) {
	health := SystemHealth{
		DatabaseStatus: false,
		Version:        version.GetVersion(),
	}

	// test database connection
	err := routes.DB.TestConnection()

	health.DatabaseStatus = err == nil
	health.DatabaseStatusMessage = "OK"

	if err != nil {
		health.DatabaseStatusMessage = err.Error()
	}

	// Build json response
	j, err := json.Marshal(health)
	if err != nil {
		routes.Logger.Error("failed to marshal health response", "error", err)
		http.Error(w, BuildResponseMessage(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)

		return
	}

	// Set json header and http code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// print json response
	_, err = w.Write(j)
	if err != nil {
		routes.Logger.Error("failed to write health response", "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}
}
