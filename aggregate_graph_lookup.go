package main

import (
	"context"
	"fmt"

	"github.com/simagix/keyhole/mdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
 * https://docs.mongodb.com/manual/reference/operator/aggregation/graphLookup/
 */
func TestAggregateGraphLookup() {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur *mongo.Cursor
	var ctx = context.Background()
	var doc bson.M

	client = getMongoClient()
	defer client.Disconnect(ctx)
	seedCarsData(client, dbName)

	pipeline := `
	[{
		"$graphLookup": {
			"from": "employees",
			"startWith": "$manager",
			"connectFromField": "manager",
			"connectToField": "_id",
			"as": "employeeHierarchy"
		}
	}]
	`

	collection = client.Database(dbName).Collection("employees")
	opts := options.Aggregate()
	if cur, err = collection.Aggregate(ctx, mdb.MongoPipeline(pipeline), opts); err != nil {
		fmt.Println(err)
	}
	defer cur.Close(ctx)
	count := int64(0)
	for cur.Next(ctx) {
		cur.Decode(&doc)
		fmt.Println("doc = ", doc)
		count++
	}

	if 0 == count {
		fmt.Println("no doc found")
	}
}
