package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	pb "github.com/martinmhan/crud-api-golang-grpc/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBAccesser is a Data Access Object that implements the dbAccesser interface
type MongoDBAccesser struct {
	connection *mongo.Client
	dbName     string
}

// ItemFields is a struct containing all fields of Item minus the ID
type ItemFields struct {
	Name string
}

type idFilter struct {
	id string
}

// Connect establishes a client connection to MongoDB and sets it as m.connection
func (m *MongoDBAccesser) Connect() error {
	godotenv.Load()
	connectionURI := os.Getenv("DB_CONNECTION_URI")
	dbName := os.Getenv("DB_NAME")

	if connectionURI == "" || dbName == "" {
		log.Fatal("Missing DB environment variables")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURI))
	if err != nil {
		log.Fatal("Error connection to MongoDB: ", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	m.connection = client
	m.dbName = dbName

	return nil
}

// InsertItem adds an item to mongodb collection
func (m *MongoDBAccesser) InsertItem(item interface{}) *pb.Item {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	log.Println("item: ", item)
	i, ok := item.(ItemFields)
	if !ok {
		log.Fatal("Invalid item")
	}

	insertResult, err := m.connection.Database(m.dbName).Collection("items").InsertOne(context.TODO(), i)
	if err != nil {
		log.Fatal("Error inserting item: ", err)
	}

	log.Println("insertResult: ", insertResult)
	log.Println("insertResult.InsertedID: ", insertResult.InsertedID)
	id := insertResult.InsertedID.(primitive.ObjectID).Hex()

	log.Println("id: ", id)

	return &pb.Item{
		ID:   id,
		Name: i.Name,
	}
}

// GetItem ...
func (m *MongoDBAccesser) GetItem(id interface{}) *pb.Item {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	i, ok := id.(string)
	if !ok {
		log.Fatal("Invalid id")
	}

	f := idFilter{
		id: i,
	}

	it := pb.Item{
		ID: i,
	}
	err := m.connection.Database(m.dbName).Collection("items").FindOne(context.TODO(), f).Decode(&it)
	if err != nil {
		log.Fatalf("Error decoding item: %v", err)
	}

	return &it
}

// UpdateItem updates an item in the mongo collection
func (m *MongoDBAccesser) UpdateItem(id interface{}, updates interface{}) error {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	idStr, ok := id.(string)
	if !ok {
		log.Fatal("Invalid id")
	}

	id, _ = primitive.ObjectIDFromHex(idStr)

	upd, ok := updates.(ItemFields)
	if !ok {
		log.Fatal("Invalid updates")
	}

	f := bson.M{"_id": id}
	u := bson.M{
		"$set": bson.M{"name": upd.Name},
	}

	log.Println("u: ", u)
	r, err := m.connection.Database(m.dbName).Collection("items").UpdateOne(context.TODO(), f, u)
	if err != nil {
		return err
	}

	log.Println("*r: ", *r)
	return nil
}

// DeleteItem deletes an item from the mongo collection by ID
func (m *MongoDBAccesser) DeleteItem(id interface{}) error {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	s, ok := id.(string)
	if !ok {
		log.Fatal("Invalid id")
	}

	i := idFilter{id: s}
	_, err := m.connection.Database(m.dbName).Collection("items").DeleteOne(context.TODO(), i)
	if err != nil {
		log.Fatalf("Error deleting item: %v", err)
		return err
	}

	return nil
}
