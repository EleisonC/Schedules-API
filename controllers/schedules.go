package controllers

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/gorilla/mux"

	"github.com/EleisonC/Schedules-API/models"
	"github.com/EleisonC/Schedules-API/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func CreateSchedule(w http.ResponseWriter, r *http.Request) {
	var scDataHolder models.SchedulesObj
	scDataHolder.ID = primitive.NewObjectID()
	scDataHolder.DelTrue = false

	err := utils.ParseBody(r, &scDataHolder)
	if err != nil { 
		utils.ErrorHandler(w, err, "Failed To Parse The Body")
		return 
	}

	ownerObjID, err := primitive.ObjectIDFromHex(scDataHolder.OwnerID)
	if err != nil { 
		utils.ErrorHandler(w, err, "Invalid Owner Id ")
		return
	}

	scDataHolder.DurTime = scDataHolder.EndTime.Sub(scDataHolder.StartTime)
	if scDataHolder.DurTime < 0 {
		err101 := errors.New("Error: Invaid Start Time or End Time")
		utils.ErrorHandler(w, err101, "There was an error validating your data")
		return
	}

	if validateErr := validate.Struct(&scDataHolder); validateErr != nil {
		utils.ErrorHandler(w, validateErr, "There was an error validating your data")
		return
	}

	// validate the owner exists
	if  errValOwnerID := scDataHolder.ValidateOwnerDB(w, ownerObjID); errValOwnerID != nil {
		utils.ErrorHandler(w, errValOwnerID, "There was an error validating your data")
		return
	}
	// validate schedule type
	scTyobjID, err := primitive.ObjectIDFromHex(scDataHolder.TypeID)
	if err != nil { 
		utils.ErrorHandler(w, err, "Invalid Schedule Type Id ")
		return
	}
	
	if  errValOwnerID := scDataHolder.ValidateScTypDB(w, scTyobjID); errValOwnerID != nil {
		utils.ErrorHandler(w, errValOwnerID, "There was an error validating your data")
		return
	}

	res := scDataHolder.AddToDB(w, &scDataHolder)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetAllSchedules(w http.ResponseWriter, r *http.Request) {
	var scDataHolder models.SchedulesObj
	params := mux.Vars(r)
	ownerID := params["ownerID"]

	objID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retrieve Data 1")
	}

	// validate the owner exists
	if  errValOwnerID := scDataHolder.ValidateOwnerDB(w, objID); errValOwnerID != nil {
		utils.ErrorHandler(w, errValOwnerID, "There was an error validating your data")
		return
	}

	result := scDataHolder.GetFromDB(w, ownerID)
	res, err := json.Marshal(result)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data 2")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetOneSchedule(w http.ResponseWriter, r *http.Request) {
	var scDataHolder models.SchedulesObj
	params := mux.Vars(r)
	ownerID := params["ownerID"]
	scID := params["scID"]

	objOwnerID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retrieve Data")
		return
	}
	if  errValOwnerID := scDataHolder.ValidateOwnerDB(w, objOwnerID); errValOwnerID != nil {
		utils.ErrorHandler(w, errValOwnerID, "There was an error validating your data")
		return
	}

	objScID, err := primitive.ObjectIDFromHex(scID)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retrieve Data")
		return
	}
	if  errScID := scDataHolder.ValidateScDB(w, objScID); errScID != nil {
		utils.ErrorHandler(w, errScID, "Schedule might not exist")
		return
	}

	result := scDataHolder.GetOneFrmDB(w, ownerID, objScID)
	res, err := json.Marshal(result)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	var scDataHolder models.SchedulesObj
	params := mux.Vars(r)
	scIDobj := params["scID"]
	scDataHolder.DelTrue = false

	err := utils.ParseBody(r, &scDataHolder)
	if err != nil { 
		utils.ErrorHandler(w, err, "Failed To Parse The Body")
		return 
	}

	scDataHolder.DurTime = scDataHolder.EndTime.Sub(scDataHolder.StartTime)
	if scDataHolder.DurTime < 0 {
		err = errors.New("Error: Invaid Start Time or End Time")
		utils.ErrorHandler(w, err, "There was an error validating your data")
		return
	}
	
	if validateErr := validate.Struct(&scDataHolder); validateErr != nil {
		utils.ErrorHandler(w, validateErr, "There was an error validating your data")
		return
	}

	ownerObjID, err := primitive.ObjectIDFromHex(scDataHolder.OwnerID)
	if err != nil {
		utils.ErrorHandler(w, err, "Invalid Owner Id ")
		return
	}

	objID, err := primitive.ObjectIDFromHex(scIDobj)
	if err != nil {
		utils.ErrorHandler(w, err, "Invalid schedule Id")
		return
	}

	scTyobjID, err := primitive.ObjectIDFromHex(scDataHolder.TypeID)
	if err != nil { 
		utils.ErrorHandler(w, err, "Invalid schedule Type Id ")
		return
	}


	// validate the owner exists
	if  errValOwnerID := scDataHolder.ValidateOwnerDB(w, ownerObjID); errValOwnerID != nil {
		utils.ErrorHandler(w, errValOwnerID, "There was an error validating your data")
		return
	}

	// validate schedule type
	if  errScTyID := scDataHolder.ValidateScTypDB(w, scTyobjID); errScTyID != nil {
		utils.ErrorHandler(w, errScTyID, "Invalid schedule type")
		return
	}

	// validate the schedule exists
	if  errScID := scDataHolder.ValidateScDB(w, objID); errScID != nil {
		utils.ErrorHandler(w, errScID, "Schedule might not exist")
		return
	}

	result := scDataHolder.UpdateScheduleFrmDB(w, objID, &scDataHolder, scDataHolder.OwnerID)

	if result.MatchedCount == 0 {
		message := "The Schedule with ID " + scIDobj + " has not been updated or does not exist"
		count := result.ModifiedCount
		res := utils.ResMessage{
			Message: message,
			Count: count,
		}
		finalRes, err := json.Marshal(res)
		if err != nil {
			utils.ErrorHandler(w, err, "Error marshalling the data")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(finalRes)
	}
	message := "The Schedule with ID " + scIDobj + " has been updated"
	count := result.ModifiedCount

	semiRes := utils.ResMessage{
		Message: message,
		Count: count,
	}

	res, err := json.Marshal(semiRes)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	var scDataHolder models.SchedulesObj
	params := mux.Vars(r)
	ownerID := params["ownerID"]
	scID := params["scID"]

	objOwnerID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retrieve Data")
		return
	}
	if  errValOwnerID := scDataHolder.ValidateOwnerDB(w, objOwnerID); errValOwnerID != nil {
		utils.ErrorHandler(w, errValOwnerID, "There was an error validating your data")
		return
	}

	objScID, err := primitive.ObjectIDFromHex(scID)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retrieve Data")
		return
	}
	if  errScID := scDataHolder.ValidateScDB(w, objScID); errScID != nil {
		utils.ErrorHandler(w, errScID, "Schedule might not exist")
		return
	}

	result := scDataHolder.DeleteScheduleFrmDB(w, objScID, ownerID)
	message := "The Schedule with ID " + scID + " has been deleted"
	count := result.ModifiedCount

	semiRes := utils.ResMessage{
		Message: message,
		Count: count,
	}

	res, err := json.Marshal(semiRes)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}