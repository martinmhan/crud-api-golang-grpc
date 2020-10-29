package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/martinmhan/crud-api-golang-grpc/utils"

	"github.com/joho/godotenv"
	"github.com/martinmhan/crud-api-golang-grpc/database"
	pb "github.com/martinmhan/crud-api-golang-grpc/proto"
	"github.com/martinmhan/crud-api-golang-grpc/rabbitmq"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAPIServer
	dba database.DBAccesser
}

func (s *server) CreateItem(ctx context.Context, in *pb.ItemFields) (*pb.SimpleMessage, error) {
	log.Println("Endpoint hit: CreateItem")

	m := rabbitmq.Message{
		Type:    "create",
		Payload: in,
	}

	rabbitmq.Produce(m)

	r := pb.SimpleMessage{
		Message: "Item creation accepted",
	}

	return &r, nil
}

func (s *server) ReadAllItems(ctx context.Context, in *pb.SimpleMessage) (*pb.Items, error) {
	log.Println("Endpoint hit: ReadAllItems")

	items := []*pb.Item{}
	records := s.dba.GetAllItems()
	for _, r := range records {
		items = append(items, &pb.Item{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return &pb.Items{
		Items: items,
	}, nil
}

func (s *server) ReadItem(ctx context.Context, in *pb.ItemID) (*pb.Item, error) {
	log.Println("Endpoint hit: ReadItem")

	it := s.dba.GetItem(in.ID)

	return &pb.Item{
		ID:   it.ID,
		Name: it.Name,
	}, nil
}

func (s *server) UpdateItem(ctx context.Context, in *pb.Item) (*pb.SimpleMessage, error) {
	log.Println("Endpoint hit: UpdateItem")

	m := rabbitmq.Message{
		Type:    "update",
		Payload: in,
	}

	rabbitmq.Produce(m)

	r := pb.SimpleMessage{
		Message: "Item update accepted",
	}

	return &r, nil
}

func (s *server) DeleteItem(ctx context.Context, in *pb.ItemID) (*pb.SimpleMessage, error) {
	log.Println("Endpoint hit: DeleteItem")

	m := rabbitmq.Message{
		Type:    "delete",
		Payload: in,
	}

	rabbitmq.Produce(m)

	r := pb.SimpleMessage{
		Message: "Item deletion accepted",
	}

	return &r, nil
}

func main() {
	godotenv.Load()

	port := ":" + os.Getenv("PORT")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	m := &database.MongoDBAccesser{}
	err = m.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAPIServer(s, &server{dba: m})
	err = s.Serve(lis)
	utils.FailOnError(err, "Failed to start producer server")
}
