package main

import (
    	"api/app/middleware"
	"api/config"
	"api/routes"

	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	var err error

	// Load relevant environment variables
	if err = config.LoadEnv(); err != nil {
		log.Fatal(err.Error())
	}
	port := config.GetEnvPort()

	// Setup router
	router := mux.NewRouter().StrictSlash(true)
	routes.InitRoutes(router)
	corsRouter := middleware.CORS(router)

	// Run server
	log.Fatal(http.ListenAndServe(port, corsRouter))
}