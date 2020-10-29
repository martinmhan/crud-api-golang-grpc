package database

// DBAccesser is the Data Access Object interface
type DBAccesser interface {
	Connect() error
	InsertItem(f ItemFields) Item
	GetAllItems() []Item
	GetItem(id string) Item
	UpdateItem(i Item) error
	DeleteItem(id string) error
}

// Item is a struct type containing all fields of an Item record
type Item struct {
	ID   string
	Name string
}

// ItemFields is a struct type containing all fields of Item excluding the ID
type ItemFields struct {
	Name string
}

// ItemID is a struct type containing just the ID
type ItemID struct {
	ID string
}
