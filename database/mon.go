package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
	"github.com/martinmhan/crud-api-golang-grpc/utils"
)

// MongoDBAccesser is a Data Access Object that implements the dbAccesser interface
type MongoDBAccesser struct {
	connection *mongo.Client
	dbName     string
}

// Connect establishes a client connection to MongoDB and sets it as m.connection
func (m *MongoDBAccesser) Connect() error {
	godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Missing DB environment variables")
	}

	connectionURI := "mongodb://" + dbHost + ":" + dbPort + "/"
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
func (m *MongoDBAccesser) InsertItem(f ItemFields) Item {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	insertResult, err := m.connection.Database(m.dbName).Collection("items").InsertOne(context.TODO(), f)
	if err != nil {
		log.Fatal("Error inserting item: ", err)
	}

	return Item{
		ID:   insertResult.InsertedID.(primitive.ObjectID).Hex(),
		Name: f.Name,
	}
}

// GetAllItems returns all items in the mongodb items collection
func (m *MongoDBAccesser) GetAllItems() []Item {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	f := bson.M{}
	cursor, err := m.connection.Database(m.dbName).Collection("items").Find(context.TODO(), f)
	utils.FailOnError(err, "Error getting all items")

	var records []bson.M
	err = cursor.All(context.TODO(), &records)
	utils.FailOnError(err, "Error reading all items")

	items := []Item{}
	for _, r := range records {
		items = append(items, Item{
			ID:   r["_id"].(primitive.ObjectID).Hex(),
			Name: r["name"].(string),
		})
	}

	return items
}

// GetItem reads an item from the mongodb collection given an ID
func (m *MongoDBAccesser) GetItem(id string) Item {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	record := bson.M{}

	_id, err := primitive.ObjectIDFromHex(id)
	utils.FailOnError(err, "Invalid id")

	f := bson.M{"_id": _id}
	result := m.connection.Database(m.dbName).Collection("items").FindOne(context.TODO(), f)
	err = result.Decode(&record)
	if err != nil {
		log.Printf("Error decoding document: %s", err)
		return Item{}
	}

	return Item{
		ID:   record["_id"].(primitive.ObjectID).Hex(),
		Name: record["name"].(string),
	}
}

// UpdateItem updates an item in the mongo collection
func (m *MongoDBAccesser) UpdateItem(item Item) error {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	id, err := primitive.ObjectIDFromHex(item.ID)
	utils.FailOnError(err, "Invalid ID")

	f := bson.M{"_id": id}
	u := bson.M{
		"$set": bson.M{"name": item.Name},
	}

	_, err = m.connection.Database(m.dbName).Collection("items").UpdateOne(context.TODO(), f, u)
	utils.FailOnError(err, "Error updating item")

	return nil
}

// DeleteItem deletes an item from the mongo collection by ID
func (m *MongoDBAccesser) DeleteItem(id string) error {
	if m.connection == nil {
		log.Fatal("DBAccesser is not connected to a database")
	}

	i, err := primitive.ObjectIDFromHex(id)
	utils.FailOnError(err, "Invalid ID")

	f := bson.M{"_id": i}
	log.Println("f: ", f)
	_, err = m.connection.Database(m.dbName).Collection("items").DeleteOne(context.TODO(), f)
	utils.FailOnError(err, "Failed to delete item")

	return nil
}
