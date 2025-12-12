package routes

import (
	"api/app/lists"
	"api/app/tasks"

	"github.com/gorilla/mux"
)

func InitRoutes(router *mux.Router) {
	// Lists
	router.HandleFunc("/api/lists", lists.GetAllListsAPI).Methods("GET")
	router.HandleFunc("/api/lists/{uuid}", lists.GetListAPI).Methods("GET")
	router.HandleFunc("/api/lists", lists.CreateListAPI).Methods("POST")
	router.HandleFunc("/api/lists/{uuid}", lists.UpdateListAPI).Methods("PUT")
	router.HandleFunc("/api/lists/{uuid}", lists.DeleteListAPI).Methods("DELETE")

	// Tasks
	router.HandleFunc("/api/tasks/{uuid}", tasks.GetTaskAPI).Methods("GET")
	router.HandleFunc("/api/tasks", tasks.CreateTaskAPI).Methods("POST")
	router.HandleFunc("/api/tasks/{uuid}", tasks.UpdateTaskAPI).Methods("PUT")
	router.HandleFunc("/api/tasks/{uuid}", tasks.DeleteTaskAPI).Methods("DELETE")
}