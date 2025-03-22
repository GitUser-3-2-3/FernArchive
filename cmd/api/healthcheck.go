package main

import (
	"net/http"
)

func (bknd *backend) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{"status": "available",
		"system_info": map[string]string{
			"version":     version,
			"environment": bknd.config.env,
		},
	}
	err := bknd.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}
