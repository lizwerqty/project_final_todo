package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		putTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func afterNow(now, date time.Time) bool {
	y1, m1, d1 := now.Date()
	y2, m2, d2 := date.Date()

	today := time.Date(y1, m1, d1, 0, 0, 0, 0, now.Location())
	taskDay := time.Date(y2, m2, d2, 0, 0, 0, 0, date.Location())

	return today.After(taskDay)
}

func checkDate(task *db.Task) error {
	var next string
	layout := "20060102"
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(layout)
	}
	t, err := time.Parse(layout, task.Date)
	if err != nil {
		return err
	}
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(layout)
		} else {
			task.Date = next
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "empty title"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "empty id"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	id := task.ID
	if id == "" {
		writeJSON(w, map[string]string{"error": "empty id"})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "empty title"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "empty id"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJSON(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "empty id"})
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, map[string]string{})
}
