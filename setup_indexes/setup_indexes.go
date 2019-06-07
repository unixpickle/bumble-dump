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

	log.Println("Removing duplicates...")
	removeDuplicates(db.Collection("profiles"))
	removeDuplicates(db.Collection("photos"))

	log.Println("Creating indices...")
	createUniqueID(db.Collection("profiles"))
	createUniqueID(db.Collection("photos"))

	log.Println("Creating location field...")
	createLocationField(db.Collection("profiles"))
	log.Println("Creating location index...")
	createLocationIndex(db.Collection("profiles"))
}

func removeDuplicates(coll *mongo.Collection) {
	log.Println("scanning table...")
	res, err := coll.Find(context.Background(), bson.D{}, nil)
	essentials.Must(err)
	ids := map[string][]interface{}{}
	for res.Next(context.Background()) {
		var obj map[string]interface{}
		if err := res.Decode(&obj); err != nil {
			log.Fatal(err)
		}
		ids[obj["id"].(string)] = append(ids[obj["id"].(string)], obj["_id"])
	}
	log.Println("removing duplicates...")
	numRemoved := 0
	for _, uids := range ids {
		if len(uids) > 1 {
			for _, uid := range uids[1:] {
				numRemoved += 1
				coll.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: uid}})
			}
		}
	}
	log.Println("removed", numRemoved, "from", coll.Name())
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

func createLocationField(coll *mongo.Collection) {
	res, err := coll.Find(context.Background(), bson.D{
		{
			Key:   "location",
			Value: bson.D{{Key: "$exists", Value: false}},
		},
	}, nil)
	essentials.Must(err)
	for res.Next(context.Background()) {
		var u bumble.User
		if err := res.Decode(&u); err != nil {
			log.Fatal(err)
		}
		u.SetLocation()
		_, err := coll.ReplaceOne(context.Background(), bson.D{{Key: "id", Value: u.ID}}, &u, nil)
		if err != nil {
			log.Fatal(err)
		}
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
