package api

import "net/http"

var todoPassword string

func Init(password string) {
	todoPassword = password

	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/signin", signInHandler)
}
