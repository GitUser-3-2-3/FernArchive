package main

import (
	"fmt"
	"net/http"
)

func (bknd *backend) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "status: available")
	_, _ = fmt.Fprintf(w, "version: %f\n", version)
	_, _ = fmt.Fprintf(w, "environment: %s\n", bknd.config.env)
}
