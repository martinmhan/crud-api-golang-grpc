package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/martinmhan/crud-api-golang-grpc/database"

	pb "github.com/martinmhan/crud-api-golang-grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedAPIServer
	dba database.DBAccesser
}

func (s *server) CreateItem(ctx context.Context, in *pb.CreateItemRequest) (*pb.Item, error) {
	log.Println("Endpoint hit: CreateItem")

	log.Println("in: ", in)
	i := database.ItemFields{
		Name: in.Name,
	}

	it := s.dba.InsertItem(i)
	log.Println("it: ", it)

	return it, nil
}

func (s *server) ReadItem(ctx context.Context, in *pb.ReadItemRequest) (*pb.Item, error) {
	log.Println("Endpoint hit: ReadItem")

	it := s.dba.GetItem(in.ID)

	return it, nil
}

func (s *server) UpdateItem(ctx context.Context, in *pb.UpdateItemRequest) (*pb.MessageResponse, error) {
	log.Println("Endpoint hit: UpdateItem")

	updates := database.ItemFields{
		Name: in.Name,
	}

	err := s.dba.UpdateItem(in.ID, updates)
	if err != nil {
		log.Fatalf("Error updating item: %v", err)
	}

	r := pb.MessageResponse{
		Message: "Item updated successfully",
	}

	return &r, nil
}

func (s *server) DeleteItem(ctx context.Context, in *pb.DeleteItemRequest) (*pb.MessageResponse, error) {
	log.Println("Endpoint hit: DeleteItem")

	err := s.dba.DeleteItem(in.ID)
	if err != nil {
		log.Fatalf("Error deleting item: %v", err)
	}

	r := pb.MessageResponse{
		Message: "Item deleted successfully",
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
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
