package storage

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoDBUri        string = "mongodb://127.0.0.1:27017"
	mongoDBName       string = "proxy"
	mongoDBCollection string = "users"
)

// Storage ...
type Db struct {
	conn           mongo.Client
	userCollection mongo.Collection
}

// Connect ...
func (s *Db) Connect() {

	ctx := context.TODO()
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBUri))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create connect
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to MongoDB!")

	s.userCollection = client.Database(mongoDBName).Collection(mongoDBCollection)
}

func (s *Db) FindUser(name string) *ProxyUser {

	s.userCollection.FinfOne(bson.M{
		"user"
	})
}
