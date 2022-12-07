package routes

import (
	"github.com/gorilla/mux"
	"github.com/EleisonC/Schedules-API/controllers"

)

var RegisterScheduleRoutes = func(router *mux.Router) {
	router.HandleFunc("/createScheduleType", controllers.CreateScheduleType).Methods("POST")
	router.HandleFunc("/getallScheduleTypes", controllers.GetAllScheduleTypes).Methods("GET")
	router.HandleFunc("/deleteScheduleTypes/{scTyID}", controllers.DeleteScheduleType).Methods("DELETE")
}
