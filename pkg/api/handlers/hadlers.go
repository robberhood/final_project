package handlers

import (
	"net/http"
	"time"

	"github.com/robberhood/final_project/pkg/api/service"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	if nowStr == "" {
		nowStr = time.Now().Format(service.DateFormat)
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(service.DateFormat, nowStr)
	if err != nil {
		http.Error(w, "field 'now' has invalid format", http.StatusBadRequest)
		return
	}

	nextDay, err := service.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDay))

}
