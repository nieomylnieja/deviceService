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
}

type DeviceDao interface {
	AddDevice(device *DevicePayload, ctx context.Context) (primitive.ObjectID, error)
	GetDevice(id string, ctx context.Context) (*Device, error)
	GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error)
	GetAllDevices(ctx context.Context) ([]Device, error)
}

func NewDao() *Dao {
	mongodbURI := os.Getenv("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Panicf("couldnt create client with the uri: %s :%+v", mongodbURI, err.Error())
	}
	collection := client.Database(os.Getenv("MONGODB_NAME")).Collection("devices")
	dao := &Dao{
		mongoClient: client,
		collection:  collection,
	}
	dao.connect(context.Background())
	return dao
}

func (db *Dao) connect(ctx context.Context) {
	err := db.mongoClient.Connect(ctx)
	if err != nil {
		log.Panicf("couldn't connect to db: %+v", err.Error())
	}
	err = db.mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Panicf("connection with db was not established properly: %+v", err.Error())
	}
}

func (db *Dao) AddDevice(device *DevicePayload, ctx context.Context) (primitive.ObjectID, error) {
	dev := Device{
		Id:       primitive.NewObjectID(),
		Name:     device.Name,
		Value:    device.Value,
		Interval: device.Interval,
	}
	result, err := db.collection.InsertOne(ctx, dev)
	if err != nil {
		log.Printf("%v was not added to db: %+v", dev, err.Error())
		return [12]byte{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (db *Dao) GetDevice(id string, ctx context.Context) (*Device, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	findResult := db.collection.FindOne(ctx, bson.M{"_id": objID})
	if err := findResult.Err(); err != nil {
		return nil, err
	}
	var dev Device
	err = findResult.Decode(&dev)
	if err != nil {
		return nil, err
	}

	return &dev, nil
}

func (db *Dao) GetAllDevices(ctx context.Context) ([]Device, error) {
	allDevices := []Device{}
	cursor, err := db.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &allDevices)
	return allDevices, err
}

func (db *Dao) GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error) {
	lower, upper := setPageBoundsToInt64(limit, page)
	paginatedDevices := []Device{}
	opts := options.FindOptions{}

	cursor, err := db.collection.Find(ctx, bson.D{},
		opts.SetSkip(lower),
		opts.SetLimit(upper))
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &paginatedDevices)
	return paginatedDevices, err
}
