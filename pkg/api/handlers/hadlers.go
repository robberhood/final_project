package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/robberhood/final_project/config"
	"github.com/robberhood/final_project/pkg/api/service"
	"github.com/robberhood/final_project/pkg/db"
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

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	var task config.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		service.WriteError(w, err)
		return
	}

	if task.Title == "" {
		service.WriteError(w, errors.New("title must be not empty"))
		return
	}

	if err := service.CheckDate(&task); err != nil {
		service.WriteError(w, err)
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		service.WriteError(w, err)
		return
	}

	service.WriteJson(w, map[string]any{"id": id})

}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		service.WriteError(w, err)
		return
	}
	service.WriteJson(w, config.TasksResp{Tasks: tasks})
}
