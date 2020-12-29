package dbs

import "go.mongodb.org/mongo-driver/mongo"

// Mongo mongodb setup
type Mongo struct {
	Schema      string `json:",omitempty"`
	MaxPoolSize uint64 `json:",omitempty"`
	MinPoolSize uint64 `json:",omitempty"`
}

var MongoDB *mongo.Client

func defaultsMongo() *Mongo {
	return &Mongo{
		Schema:      "mongodb://127.0.0.1:27017/summer",
		MaxPoolSize: 20,
		MinPoolSize: 5,
	}
}
