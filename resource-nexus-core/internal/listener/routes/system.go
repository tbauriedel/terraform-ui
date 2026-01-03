package routes

import (
	"encoding/json"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/version"
)

type SystemHealth struct {
	Database bool   `json:"databaseConnection"`
	Version  string `json:"version"`
}

func (routes *Routes) Health(w http.ResponseWriter, r *http.Request) {
	health := SystemHealth{
		Database: false,
		Version:  version.GetVersion(),
	}

	// test database connection
	err := routes.DB.TestConnection()
	if err == nil {
		health.Database = true
	}

	// Build json response
	j, err := json.Marshal(health)
	if err != nil {
		panic(err)
	}

	// Set json header and http code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// print json response
	_, err = w.Write(j)
	if err != nil {
		panic(err)
	}
}
