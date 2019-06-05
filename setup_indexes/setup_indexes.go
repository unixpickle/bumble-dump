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
}

func removeDuplicates(coll *mongo.Collection) {
	res, err := coll.Find(context.Background(), bson.D{}, nil)
	essentials.Must(err)
	ids := map[string][]interface{}{}
	for {
		var obj struct {
			RealID interface{} `bson:"_id"`
			ID     string      `bson:"id"`
		}
		if err := res.Decode(&obj); err != nil {
			if err == mongo.ErrNoDocuments {
				break
			}
			ids[obj.ID] = append(ids[obj.ID], obj.RealID)
		}
	}
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
	coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}, nil)
}
