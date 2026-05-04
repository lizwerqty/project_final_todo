package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

const limit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")
	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
