package routes

import (
	"api/app/lists"
	"api/app/tasks"

	"github.com/gorilla/mux"
)

func InitRoutes(router *mux.Router) {
	// Lists
	router.HandleFunc("/lists", lists.GetAllListsAPI).Methods("GET")
	router.HandleFunc("/lists/{uuid}", lists.GetListAPI).Methods("GET")
	router.HandleFunc("/lists", lists.CreateListAPI).Methods("POST")
	router.HandleFunc("/lists", lists.UpdateListsAPI).Methods("PUT")
	router.HandleFunc("/lists/{uuid}", lists.DeleteListAPI).Methods("DELETE")

	// Tasks
	router.HandleFunc("/tasks/{uuid}", tasks.GetTaskAPI).Methods("GET")
	router.HandleFunc("/tasks", tasks.CreateTaskAPI).Methods("POST")
	router.HandleFunc("/tasks", tasks.UpdateTasksAPI).Methods("PUT")
	router.HandleFunc("/tasks/{uuid}", tasks.DeleteTaskAPI).Methods("DELETE")
}