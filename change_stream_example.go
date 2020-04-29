package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

var collection = "oplogs"

// example: argos "mongodb://localhost:27017/argos?replicaSet=replset" students '[{"$match": {"operationType": "update"}}]'
func silent(doc bson.M) {
	fmt.Println("------------I feel something--------------")
	fmt.Println("doc = ", doc)
}

func ChangeStreamClient() {
	var err error
	var client *mongo.Client
	var cs connstring.ConnString
	var ctx = context.Background()
	// var uri = "mongodb://localhost:27017/argos?replicaSet=replset"
	if os.Getenv("DATABASE_URL") != "" {
		uri = os.Getenv("DATABASE_URL")
	}
	if cs, err = connstring.Parse(uri); err != nil {
		fmt.Println(err)
	}
	client = GetMongoClient()
	defer client.Disconnect(ctx)
	var pipeline []bson.D
	pipeline = mongo.Pipeline{}
	c := client.Database(cs.Database).Collection(collection)
	c.InsertOne(ctx, bson.M{"city": "Atlanta"})

	go func(c *mongo.Collection) {
		execute(c)
		client.Disconnect(context.Background())
	}(c)

	stream := NewChangeStream()
	stream.SetPipeline(pipeline)
	// if testing.Short() {
	// 	// t.Skip("test changes stream")
	// 	stream.Watch(client, silent)
	// }
	stream.Watch(client, silent)
}

func ChangeStreamDatabase() {
	var err error
	var client *mongo.Client
	var cs connstring.ConnString
	var ctx = context.Background()
	// var uri = "mongodb://localhost:27017/argos?replicaSet=replset"
	if os.Getenv("DATABASE_URL") != "" {
		uri = os.Getenv("DATABASE_URL")
	}
	if cs, err = connstring.Parse(uri); err != nil {
		fmt.Println(err)
	}
	client = GetMongoClient()
	var pipeline []bson.D
	pipeline = mongo.Pipeline{}
	c := client.Database(cs.Database).Collection(collection)
	c.InsertOne(ctx, bson.M{"city": "Atlanta"})

	go func(c *mongo.Collection) {
		execute(c)
		client.Database(cs.Database).Drop(context.Background())
	}(c)

	stream := NewChangeStream()
	stream.SetDatabase(cs.Database)
	stream.SetPipeline(pipeline)
	// if testing.Short() {
	// 	t.Skip("test changes stream")
	// 	stream.Watch(client, silent)
	// }
	stream.Watch(client, silent)
}

func ChangeStreamCollection() {
	var err error
	var client *mongo.Client
	var cs connstring.ConnString
	var ctx = context.Background()
	var uri = "mongodb://localhost:27017/argos?replicaSet=replset"
	if os.Getenv("DATABASE_URL") != "" {
		uri = os.Getenv("DATABASE_URL")
	}
	if cs, err = connstring.Parse(uri); err != nil {
		fmt.Println(err)
	}
	fmt.Println("cs = ", cs)
	fmt.Println("cs.Database = ", cs.Database)

	client = GetMongoClient()
	var pipeline []bson.D
	pipeline = mongo.Pipeline{}
	c := client.Database(cs.Database).Collection(collection)
	c.InsertOne(ctx, bson.M{"city": "Atlanta"})

	go func(c *mongo.Collection) {
		execute(c)
	}(c)

	stream := NewChangeStream()
	stream.SetCollection(collection)
	stream.SetDatabase(cs.Database)
	stream.SetPipeline(pipeline)
	// if testing.Short() {
	// 	// t.Skip("test changes stream")
	// 	stream.Watch(client, silent)
	// }
	stream.Watch(client, silent)
}

func ChangeStreamCollectionWithPipeline() {
	var err error
	var client *mongo.Client
	var cs connstring.ConnString
	var ctx = context.Background()
	// var uri = "mongodb://localhost:27017/argos?replicaSet=replset"
	if os.Getenv("DATABASE_URL") != "" {
		uri = os.Getenv("DATABASE_URL")
	}
	if cs, err = connstring.Parse(uri); err != nil {
		fmt.Println(err)
	}
	client = GetMongoClient()
	// var pipeline = mdb.MongoPipeline(`[{"$match": {"operationType": {"$in": ["update", "delete"] } }}]`)
	var pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
		bson.A{
			bson.D{{"operationType", "insert"}}}}},
	}}}
	c := client.Database(cs.Database).Collection(collection)
	c.InsertOne(ctx, bson.M{"city": "Atlanta"})

	go func(c *mongo.Collection) {
		execute(c)
	}(c)
	stream := NewChangeStream()

	stream.SetCollection(collection)
	stream.SetDatabase(cs.Database)
	stream.SetPipeline(pipeline)
	// if testing.Short() {
	// 	t.Skip("test changes stream")
	// 	stream.Watch(client, silent)
	// }
	stream.Watch(client, silent)

}

func execute(c *mongo.Collection) {
	time.Sleep(4 * time.Second) // wait for change stream to init
	var doc = bson.M{"_id": primitive.NewObjectID(), "hometown": "Atlanta"}
	c.InsertOne(context.Background(), doc)
	var update bson.M
	json.Unmarshal([]byte(`{ "$set": {"year": 1998}}`), &update)
	c.UpdateOne(context.Background(), bson.M{"_id": doc["_id"]}, update)
	c.DeleteMany(context.Background(), bson.M{"hometown": "Atlanta"})
	time.Sleep(1 * time.Second) // wait for CS to print messages
	c.Drop(context.Background())
}