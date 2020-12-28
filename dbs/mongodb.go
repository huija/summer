package dbs

import "go.mongodb.org/mongo-driver/mongo"

// Mongo mongodb setup
type Mongo struct {
	Schema      string
	MaxPoolSize uint64
	MinPoolSize uint64
}

var MongoDB *mongo.Client
