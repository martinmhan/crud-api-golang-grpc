package database

import pb "github.com/martinmhan/crud-api-golang-grpc/proto"

// DBAccesser is the Data Access Object interface
type DBAccesser interface {
	Connect() error
	InsertItem(item interface{}) *pb.Item
	GetItem(id interface{}) *pb.Item
	UpdateItem(id interface{}, updates interface{}) error
	DeleteItem(id interface{}) error
}
