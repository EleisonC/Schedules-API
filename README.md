## MyDogo Schedules API
This is a Golang Gorilla-Mux API of MyDogo that handles:
1. Scheduling
2. Creating, reading and deleting of schedule types
3. Creating, reading, updating and deleting of schedules

#### To test the application and get it running, do the following:
1. Prerequisite
- Golang
- MongoDB
- A valid userID created from the MyDogo User API
 
2. Install the requirements file for all the dependencies of the application
```
$ go mod tidy
```
3. Run the application
```
$ go run main.go 
```

#### Features
Endpoint | Functionality
------------ | -------------
POST /createScheduleType | Creates a new schedule type
GET /getallScheduleTypes | Get all schedule types
DELETE /deleteScheduleTypes/{scTyID} | Delete a schedule type
POST /createSchedule | Creates an new schedlue valid ownerId and schedule type required
PUT /updateschedule/{scID} | Update a Schdeule
GET /getallschedules/{ownerID} | Get schedules for a user
GET /getschedule/{ownerID}/{scID} | Get a schedule for a user
DELETE /deleteschedule/{ownerID}/{scID} | Delete a schedule 
