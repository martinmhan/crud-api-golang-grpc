package main

import (
	"encoding/json"
	"log"

	"github.com/martinmhan/crud-api-golang-grpc/database"
	"github.com/martinmhan/crud-api-golang-grpc/rabbitmq"
	"github.com/martinmhan/crud-api-golang-grpc/utils"
)

func createItem(d database.DBAccesser, payload []byte) {
	var r database.ItemFields

	err := json.Unmarshal(payload, &r)
	utils.FailOnError(err, "Failed to parse message payload")

	d.InsertItem(r)
}

func updateItem(d database.DBAccesser, payload []byte) {
	var r database.Item

	err := json.Unmarshal(payload, &r)
	utils.FailOnError(err, "Failed to parse message payload")

	d.UpdateItem(r)
}

func deleteItem(d database.DBAccesser, payload []byte) {
	var r database.ItemID

	err := json.Unmarshal(payload, &r)
	utils.FailOnError(err, "Failed to parse message payload")

	d.DeleteItem(r.ID)
}

func main() {
	conn := rabbitmq.Connect()
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"crud",
		false,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	db := &database.MongoDBAccesser{}
	err = db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Println("Received a message")
			log.Printf("Message Type: %s", d.Type)
			log.Printf("Message Body: %s", d.Body)

			switch d.Type {
			case "create":
				createItem(db, d.Body)
			case "update":
				updateItem(db, d.Body)
			case "delete":
				deleteItem(db, d.Body)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
