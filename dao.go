package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

type Dao struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
	ctx         context.Context
}

type DeviceDao interface {
	AddDevice(device *DevicePayload) (primitive.ObjectID, error)
	GetDevice(id string) (*Device, error)
	GetPaginatedDevices(limit int, page int) ([]Device, error)
	GetAllDevices() ([]Device, error)
}

func NewDao() *Dao {
	mongodbURI := os.Getenv("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Panicf("couldnt create client with the uri: %s :%+v", mongodbURI, err.Error())
	}
	collection := client.Database(os.Getenv("MONGODB_NAME")).Collection("devices")
	return &Dao{
		mongoClient: client,
		collection:  collection,
		ctx:         context.Background(),
	}
}

func (db *Dao) ConnectToDB() {
	err := db.mongoClient.Connect(db.ctx)
	if err != nil {
		log.Panicf("couldn't connect to db: %+v", err.Error())
	}
	err = db.mongoClient.Ping(db.ctx, readpref.Primary())
	if err != nil {
		log.Panicf("connection with db was not established properly: %+v", err.Error())
	}
}

func (db *Dao) AddDevice(device *DevicePayload) (primitive.ObjectID, error) {
	dev := Device{
		Id:       primitive.NewObjectID(),
		Name:     device.Name,
		Value:    device.Value,
		Interval: device.Interval,
	}
	result, err := db.collection.InsertOne(db.ctx, dev)
	if err != nil {
		log.Printf("%v was not added to db: %+v", dev, err.Error())
		return [12]byte{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (db *Dao) GetDevice(id string) (*Device, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	findResult := db.collection.FindOne(db.ctx, bson.M{"_id": objID})
	if err := findResult.Err(); err != nil {
		return nil, err
	}
	dev := Device{}
	err = findResult.Decode(dev)
	if err != nil {
		return nil, err
	}
	return &dev, nil
}

func (db *Dao) GetAllDevices() ([]Device, error) {
	allDevices := []Device{}
	cursor, err := db.collection.Find(db.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(db.ctx, allDevices)
	return allDevices, err
}

func (db *Dao) GetPaginatedDevices(limit, page int) ([]Device, error) {
	count, err := db.collection.CountDocuments(db.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	lower, upper := setPageBounds(int64(limit), int64(page), count)
	paginatedDevices := []Device{}
	opts := options.FindOptions{}
	cursor, err := db.collection.Find(db.ctx, bson.D{}, opts.SetSkip(lower), opts.SetLimit(upper))
	if err != nil {
		return nil, err
	}
	err = cursor.All(db.ctx, paginatedDevices)
	return paginatedDevices, err
}
