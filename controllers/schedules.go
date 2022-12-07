package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/EleisonC/Schedules-API/models"
	"github.com/EleisonC/Schedules-API/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	scDataHolder models.SchedulesObj
)
var validate = validator.New()

func CreateScheduleType(w http.ResponseWriter, r *http.Request) {
	// parse data
	var scTyDataHolder models.SchedulesType
	scTyDataHolder.ID = primitive.NewObjectID()
	err := utils.ParseBody(r, &scTyDataHolder)
	if err != nil { 
		utils.ErrorHandler(w, err, "Failed To Parse The Body")
		return 
	}

	//validate data
	if validateErr := validate.Struct(&scTyDataHolder); validateErr != nil {
		utils.ErrorHandler(w, validateErr, "There was an error validating your data")
		return
	}

	// Add to database
	res := scTyDataHolder.AddTypeToDB(w, &scTyDataHolder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetAllScheduleTypes(w http.ResponseWriter, r *http.Request) {
	var scTyDataHolder models.SchedulesType
	result := scTyDataHolder.GetAllScheduleTyps(w)
	res, err := json.Marshal(result)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteScheduleType(w http.ResponseWriter, r *http.Request) {
	var scTyDataHolder models.SchedulesType
	params := mux.Vars(r)
	scTyID := params["scTyID"]

	objId, err := primitive.ObjectIDFromHex(scTyID)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Delete Data")
		return
	}
	result := scTyDataHolder.DeleteScheduleTypeDB(w, objId)

	if result.DeletedCount == 0 {
		message := "The Schedule Type with ID " + scTyID + " has not been deleted from the DB or does not exist"
		count := result.DeletedCount
		res := utils.ResMessage{
			Message: message,
			Count: count,
		}
		finalRes, err := json.Marshal(res)
		if err != nil {
			utils.ErrorHandler(w, err, "Error Creating Response")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(finalRes)
		return
	}

	message := "The Schedule Type with ID " + scTyID + " has not been deleted from the DB or does not exist"
		count := result.DeletedCount
		res := utils.ResMessage{
			Message: message,
			Count: count,
		}
	finalRes, err := json.Marshal(res)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Creating Response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(finalRes)
}