package models

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/EleisonC/Schedules-API/configs"
	"github.com/EleisonC/Schedules-API/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SchedulesType struct {
	ID primitive.ObjectID `bson:"_id"`
	Name string `bson:"name" validate:"required"`
	Description string `bson:"description" validate:"required"`
}

var scheduleType SchedulesType

type SchedulesObj struct {
	ID primitive.ObjectID `bson:"_id"`
	TypeID string `bson:"typeId" validate:"required"`
	OwnerID string  `bson:"ownerId" validate:"required"`
	AttendeesID []string `bson:"attendeeIds"`
	DogsID []string `bson:"dogIds"`
	Name string `bson:"name" validate:"required"`
	Description string `bson:"description" validate:"required"`
	StartTime time.Time `bson:"startTime" validate:"required"`
	EndTime time.Time `bson:"endTime" validate:"required"`
	DurTime time.Duration `bson:"durTime"`
	DelTrue bool `bson:"delValue, omitempty"`
}

var SchedulesDoc *mongo.Collection = configs.GetCollection(configs.DB, "SchedulesDoc")
var SchedulesTypeDoc *mongo.Collection = configs.GetCollection(configs.DB, "SchedulesType")
var dogOwnerDoc *mongo.Collection = configs.GetCollection(configs.DB, "DogOwner")

func (s SchedulesObj) AddToDB(w http.ResponseWriter, seObj *SchedulesObj)[]byte {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := SchedulesDoc.InsertOne(ctx, *seObj)
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

func (s SchedulesObj) GetFromDB(w http.ResponseWriter, ownerObjID string)[]SchedulesObj {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"delValue": false, "ownerId": ownerObjID}
	results, err := SchedulesDoc.Find(ctx, filter)
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

func (s SchedulesObj) GetOneFrmDB(w http.ResponseWriter, owenrObjID string, scID primitive.ObjectID)SchedulesObj {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var scObj SchedulesObj
	defer cancel()
	err := SchedulesDoc.FindOne(ctx, bson.M{"_id": scID,"ownerId": owenrObjID, "delValue": false}).Decode(&scObj)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Retriving Data")
	}
	return scObj
}

func (s SchedulesObj) UpdateScheduleFrmDB(w http.ResponseWriter, scID primitive.ObjectID, scData *SchedulesObj, owenrObjID string)*mongo.UpdateResult{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scID, "delValue": false, "ownerId": owenrObjID}

	update := bson.D{{Key:"$set", Value: bson.D{{Key: "typeId", Value: scData.TypeID},
		{Key: "ownerId", Value: scData.OwnerID}, {Key: "attendeesIds", Value: scData.AttendeesID},
		{Key: "dogsId", Value: scData.DogsID}, {Key: "name", Value: scData.Name},  {Key: "description", Value: scData.Description},
		{Key: "startTime", Value: scData.StartTime}, {Key: "endTime", Value: scData.EndTime}, {Key: "durTime", Value: scData.DurTime},
		{Key: "delValue", Value: scData.DelTrue}}}}

	result, err := SchedulesDoc.UpdateOne(ctx, filter, update)
	if err != nil {
		utils.ErrorHandler(w, err, "Error Updating Data")
	}
	return result
}

func (s SchedulesObj) DeleteScheduleFrmDB(w http.ResponseWriter, scID primitive.ObjectID, ownerObjID string)*mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scID, "ownerId": ownerObjID}
	update := bson.D{{Key:"$set",  Value: bson.D{{Key: "delValue", Value: true}}}}
	delResult, err := SchedulesDoc.UpdateOne(ctx, filter, update)
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

func (s SchedulesObj) ValidateOwnerDB(w http.ResponseWriter, ownerID primitive.ObjectID) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": ownerID}
	ownerIdCount, _ := dogOwnerDoc.CountDocuments(ctx, filter)
	if ownerIdCount == 0 {
		err1 := errors.New("Error: This owner does not exist")
		utils.ErrorHandler(w, err1, "Error UnKnown User")
		return err1
	}
	return nil
}

func (s SchedulesObj) ValidateScDB(w http.ResponseWriter, scID primitive.ObjectID) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scID, "delValue": false}
	scIdCount, _ := SchedulesDoc.CountDocuments(ctx, filter)
	if scIdCount == 0 {
		err1 := errors.New("Error: This document does not exist")
		return err1
	}
	return nil
}

func (s SchedulesObj) ValidateScTypDB(w http.ResponseWriter, scTypID primitive.ObjectID) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": scTypID}
	scTypIdCount, _ := SchedulesTypeDoc.CountDocuments(ctx, filter)
	if scTypIdCount == 0 {
		err1 := errors.New("Error: This schedule type does not exist")
		return err1
	}
	return nil
}

