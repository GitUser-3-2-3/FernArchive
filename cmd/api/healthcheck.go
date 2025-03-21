package main

import (
	"net/http"
)

func (bknd *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{"status": "available",
		"environment": bknd.config.env,
		"version":     version,
	}
	err := bknd.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		bknd.logger.Error(err.Error())
		http.Error(w, "could not process your request", http.StatusInternalServerError)
	}
}
