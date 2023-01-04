package routes

import (
	"github.com/gorilla/mux"
	"github.com/EleisonC/Schedules-API/controllers"

)

var RegisterScheduleRoutes = func(router *mux.Router) {
	router.HandleFunc("/createScheduleType", controllers.CreateScheduleType).Methods("POST")
	router.HandleFunc("/getallScheduleTypes", controllers.GetAllScheduleTypes).Methods("GET")
	router.HandleFunc("/deleteScheduleTypes/{scTyID}", controllers.DeleteScheduleType).Methods("DELETE")
	router.HandleFunc("/createSchedule", controllers.CreateSchedule).Methods("POST")
	router.HandleFunc("/updateschedule/{scID}", controllers.UpdateSchedule).Methods("PUT")
	router.HandleFunc("/getallschedules/{ownerID}", controllers.GetAllSchedules).Methods("GET")
	router.HandleFunc("/getschedule/{ownerID}/{scID}", controllers.GetOneSchedule).Methods("GET")
	router.HandleFunc("/deleteschedule/{ownerID}/{scID}", controllers.DeleteSchedule).Methods("DELETE")
}
