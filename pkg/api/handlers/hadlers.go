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
	query := r.URL.Query().Get("search")
	if query != "" {
		date, err := time.Parse("02.01.2006", query)
		if err == nil {
			tasks, err := db.TasksByDate(date.Format(service.DateFormat))
			if err != nil {
				service.WriteJson(w, map[string]string{"error": err.Error()})
				return
			}
			service.WriteJson(w, config.TasksResp{Tasks: tasks})
			return
		}
		tasks, err := db.TasksByText(query)
		if err != nil {
			service.WriteJson(w, map[string]string{"error": err.Error()})
			return
		}
		service.WriteJson(w, config.TasksResp{Tasks: tasks})
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		service.WriteError(w, err)
		return
	}
	service.WriteJson(w, config.TasksResp{Tasks: tasks})
}

func TaskGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		service.WriteJson(w, map[string]string{"error": "No identifier specified"})
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		service.WriteJson(w, map[string]string{"error": "Task not found"})
		return
	}

	service.WriteJson(w, t)
}
func TaskUPDHandler(w http.ResponseWriter, r *http.Request) {
	var t config.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		service.WriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	if t.ID == "" {
		service.WriteJson(w, map[string]string{"error": "No identifier specified"})
		return
	}

	if t.Title == "" {
		service.WriteJson(w, map[string]string{"error": "title must be not empty"})
		return
	}

	if err := service.CheckDate(&t); err != nil {
		service.WriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	err := db.UpdateTask(&t)
	if err != nil {
		service.WriteJson(w, map[string]string{"error": "Task not found"})
		return
	}

	service.WriteJson(w, map[string]any{})
}

func TaskCompleteHandler(w http.ResponseWriter, r *http.Request) {
	var t *config.Task
	id := r.URL.Query().Get("id")
	if id == "" {
		service.WriteJson(w, map[string]string{"error": "No identifier specified"})
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		service.WriteJson(w, map[string]string{"error": "Task not found"})
		return
	}

	if t.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			service.WriteJson(w, map[string]string{"error": "Task not found"})
			return
		}
	} else {
		taskDate, err := time.Parse(service.DateFormat, t.Date)
		if err != nil {
			service.WriteJson(w, map[string]string{"error": "Invalid task date"})
			return
		}
		next, err := service.NextDate(taskDate, t.Date, t.Repeat)
		if err != nil {
			service.WriteJson(w, map[string]string{"error": err.Error()})
			return
		}
		t.Date = next
		err = db.UpdateDate(id, next)
		if err != nil {
			service.WriteJson(w, map[string]string{"error": "Task not found"})
			return
		}
	}
	service.WriteJson(w, map[string]any{})

}

func TaskDELHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		service.WriteJson(w, map[string]string{"error": "No identifier specified"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		service.WriteJson(w, map[string]string{"error": "Task not found"})
		return
	}

	service.WriteJson(w, map[string]any{})
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var pass config.Pass

	if err := json.NewDecoder(r.Body).Decode(&pass); err != nil {
		service.WriteJson(w, map[string]string{"error": err.Error()})
		return
	}
	if pass.Password != config.Password {
		service.WriteJson(w, map[string]string{"error": "Incorrect password"})
		return
	}
	pass.Token = service.GenerateToken()
	service.WriteJson(w, pass)
}
