package models

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EleisonC/Schedules-API/configs"
	"github.com/EleisonC/Schedules-API/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type SchedulesType struct {
	ID primitive.ObjectID `bson:"_id"`
	Name string `bson:"name" validate:"required"`
	Description string `bson:"description" validate:"required"`
}

var scheduleType SchedulesType

type SchedulesObj struct {
	ID primitive.ObjectID `bson:"_id"`
	TypeID primitive.ObjectID `bson:"type" validate:"required"`
	OwnerID primitive.ObjectID  `bson:"ownerId" validate:"required"`
	AttendeesID []primitive.ObjectID `bson:"attendeeIds"`
	DogsID []primitive.ObjectID `bson:"dogIds"`
	Name string `bson:"name" validate:"required"`
	Description string `bson:"description" validate:"required"`
	StartTime time.Time `bson:"startTime" validate:"required"`
	StartDate time.Time `bson:"startDate" validate:"required"`
	EndDate time.Time `bson:"endDate" validate:"required"`
	EndTime time.Time `bson:"endTime" validate:"required"`
	RemTime time.Time `bson:"remTime" validate:"required"`
	DeleteTrue string `bson:"deleteTrue, omitempty"`
}

var SchedulesDoc *mongo.Collection = configs.GetCollection(configs.DB, "SchedulesDoc")
var SchedulesTypeDoc *mongo.Collection = configs.GetCollection(configs.DB, "SchedulesType")

func (s SchedulesObj) AddToDB(w http.ResponseWriter, seObj *SchedulesObj)[]byte {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := SchedulesDoc.InsertOne(ctx, seObj)
	if err != nil {
		utils.ErrorHandler(w, err, "There was an error entering the data")
	}
	res, err:=json.Marshal(result)
	if err != nil {
		utils.ErrorHandler(w, err, "There was an error marshalling the data")
	}
	return res
}

func (st SchedulesType) AddTypeToDB(w http.ResponseWriter, scTyData *SchedulesType)[]byte {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := SchedulesTypeDoc.InsertOne(ctx, *scTyData)
	
	if err != nil {
		utils.ErrorHandler(w, err, "There was an error entering the data")
	}
	res, err:=json.Marshal(result)
	if err != nil {
		utils.ErrorHandler(w, err, "There was an error marshalling the data")
	}
	return res
}

func (s SchedulesObj) GetFromDB(w http.ResponseWriter)[]SchedulesObj {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results, err := SchedulesDoc.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
	}
	defer results.Close(ctx)
	scObjArray := []SchedulesObj{}
	for results.Next(ctx) {
		var scObj SchedulesObj
		if err := results.Decode(&scObj); err != nil {
			utils.ErrorHandler(w, err, "Error Retriving Data")
		}
		scObjArray = append(scObjArray, scObj)
	}
	return scObjArray
}

func (st SchedulesType) GetAllScheduleTyps(w http.ResponseWriter)[]SchedulesType {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results, err := SchedulesTypeDoc.Find(ctx, bson.M{})
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
	}
	defer results.Close(ctx)
	scTypObjArray := []SchedulesType{}
	for results.Next(ctx) {
		var scTyObj SchedulesType
		if err := results.Decode(&scTyObj); err != nil {
			utils.ErrorHandler(w, err, "Error Retriving Data")
		}
		scTypObjArray = append(scTypObjArray, scTyObj)
	}
	return scTypObjArray
}

func (s SchedulesObj) GetOneFrmDB(w http.ResponseWriter, owenrObjID string, scID string)SchedulesObj {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var scObj SchedulesObj
	defer cancel()
	err := SchedulesDoc.FindOne(ctx, bson.M{"_id": scID,"ownerId": owenrObjID}).Decode(&scObj)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
	}
	return scObj
}

func (s SchedulesObj) updateScheduleFrmDB(w http.ResponseWriter, scID primitive.ObjectID, scData *SchedulesObj)*mongo.UpdateResult{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scID}
	result, err := SchedulesDoc.UpdateOne(ctx, filter, scData)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Updating Data")
	}
	return result
}

func (s SchedulesObj) deleteScheduleFrmDB(w http.ResponseWriter, scID primitive.ObjectID, ownerObjID primitive.ObjectID)*mongo.DeleteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scID, "ownerId": ownerObjID}
	delResult, err := SchedulesDoc.DeleteOne(ctx, filter)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Deleting Data")
	}
	return delResult
}

func (st SchedulesType) DeleteScheduleTypeDB(w http.ResponseWriter, scID primitive.ObjectID) *mongo.DeleteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scID}
	delResult, err := SchedulesTypeDoc.DeleteOne(ctx, filter)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Deleting Data")
	}
	return delResult
}
