package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"time"
)

type Device struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Value float32 `json:"value,omitempty" bson:"value,omitempty"`
	Interval float32 `json:"interval,omitempty" bson:"interval,omitempty"`
}

var client *mongo.Client

func CreateDeviceEndpoint(response http.ResponseWriter,
						request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var device Device
	_ = json.NewDecoder(request.Body).Decode(&device)
	collection := client.Database("deviceservice").Collection("device")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, err := collection.InsertOne(ctx, device)
	if err != nil { log.Fatal(err) }
	json.NewEncoder(response).Encode(result)
}

func main() {
	fmt.Println("Starting application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("Couldn't connect!")
	}
	router := mux.NewRouter()
	router.HandleFunc("/device", CreateDeviceEndpoint).Methods("POST")
	fmt.Println("Serving...")
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}