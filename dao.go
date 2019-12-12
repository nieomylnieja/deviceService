package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"regexp"
	"runtime"
)

type Dao struct {
	mongoClient *mongo.Client
	collection  *mongo.Collection
}

type DeviceDao interface {
	AddDevice(device *DevicePayload, ctx context.Context) (primitive.ObjectID, error)
	GetDevice(id primitive.ObjectID, ctx context.Context) (*Device, error)
	GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error)
	GetAllDevices(ctx context.Context) ([]Device, error)
}

func NewDao() *Dao {
	mongodbURI := os.Getenv("MONGODB_URI")
	mongodbNAME := os.Getenv("MONGODB_NAME")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Panicf("couldnt create client with the uri: %s :%+v", mongodbURI, err.Error())
	}
	if err = verifyMongoDBName(mongodbNAME); err != nil {
		log.Panicf("incorrect name: %s: %+v", mongodbNAME, err.Error())
	}
	collection := client.Database(mongodbNAME).Collection("devices")
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

func (db *Dao) GetDevice(id primitive.ObjectID, ctx context.Context) (*Device, error) {
	findResult := db.collection.FindOne(ctx, bson.M{"_id": id})
	if err := findResult.Err(); err != nil {
		return nil, err
	}
	var dev Device
	if err := findResult.Decode(&dev); err != nil {
		return nil, err
	}
	return &dev, nil
}

func (db *Dao) GetAllDevices(ctx context.Context) ([]Device, error) {
	allDevices := make([]Device, 0)
	cursor, err := db.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &allDevices)
	return allDevices, err
}

func (db *Dao) GetPaginatedDevices(limit, page int, ctx context.Context) ([]Device, error) {
	lower, upper := setPageBoundsToInt64(limit, page)
	paginatedDevices := make([]Device, 0)
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

func verifyMongoDBName(dbName string) error {
	if !(len(dbName) < 64 && 0 < len(dbName)) {
		return errors.New("db name must not be empty")
	}
	var notAllowed string
	if runtime.GOOS == "windows" {
		notAllowed = `/\\. "$*<>:|?`
	} else {
		notAllowed = `/\\. "$`
	}
	re := regexp.MustCompile(fmt.Sprintf(`[%s]`, notAllowed))
	if re.MatchString(dbName) {
		return fmt.Errorf("db name can't contain any of these characters: %s", notAllowed)
	}
	return nil
}
