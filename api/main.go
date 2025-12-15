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
	// Setup router
	router := mux.NewRouter().StrictSlash(true)
	routes.InitRoutes(router)
	corsRouter := middleware.CORS(router)

	// Run server
	port := config.GetEnvApiPort()
	log.Fatal(http.ListenAndServe(port, corsRouter))
}