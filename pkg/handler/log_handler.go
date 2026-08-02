package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ulbithebest/BE-pendaftaran/pkg/applog"
)

func GetAppLogsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applog.List())
}
