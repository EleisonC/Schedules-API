package main

import (
	"log"
	"net/http"

	"github.com/EleisonC/Schedules-API/routes"
	"github.com/EleisonC/Schedules-API/configs"
	"github.com/gorilla/mux"
)


func main() {
	r := mux.NewRouter()
	routes.RegisterScheduleRoutes(r)
	http.Handle("/", r)
	configs.ConnectDB()
	log.Fatal(http.ListenAndServe(":8080", r))
}