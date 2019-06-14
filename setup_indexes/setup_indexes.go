// Command setup_indexes creates indexes on the MongoDB
// database.
package main

import (
	"context"
	"log"
	"time"

	"github.com/unixpickle/bumble-dump"
	"github.com/unixpickle/essentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config := bumble.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DatabaseURI))
	essentials.Must(err)
	db := client.Database("bumble")

	log.Println("Creating indices...")
	createUniqueID(db.Collection("profiles"))
	createUniqueID(db.Collection("photos"))
	createLocationIndex(db.Collection("profiles"))
}

func createUniqueID(coll *mongo.Collection) {
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createLocationIndex(coll *mongo.Collection) {
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "location", Value: 1}},
	})
	if err != nil {
		log.Fatal(err)
	}
}
