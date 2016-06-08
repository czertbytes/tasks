package main

import (
	"log"
	"net/http"

	"github.com/czertbytes/tasks"
)

func main() {
	taskMemoryStorage := tasks.NewTaskMemoryStorage()
	taskService := tasks.NewTaskStorageService(taskMemoryStorage)
	tasksHandler := tasks.NewTasksHandler(taskService)
	taskHandler := tasks.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.Handle("/tasks", tasksHandler)
	mux.Handle("/tasks/", taskHandler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
