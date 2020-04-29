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
 * count vehicles by style and display dealer name, all brands and a total count of each style
 */
func TestAggregateLookup() {
	var err error
	var client *mongo.Client
	var collection *mongo.Collection
	var cur *mongo.Cursor
	var ctx = context.Background()
	var doc bson.M

	client = getMongoClient()
	defer client.Disconnect(ctx)
	total := seedCarsData(client, dbName)

	pipeline := `
	[{
		"$group": {
			"_id": {
				"dealer": "$dealer",
				"style": "$style"
			},
			"brand": {
				"$addToSet": "$brand"
			},
			"count": {
				"$sum": 1
			}
		}
	}, {
		"$lookup": {
			"as": "dealer",
			"from": "dealers",
			"let": {
				"dealerId": "$_id.dealer"
			},
			"pipeline": [{
				"$match": {
					"$expr": {
						"$eq": [
							"$_id", "$$dealerId"
						]
					}
				}
			}, {
				"$project": {
					"_id": 0,
					"name": 1
				}
			}]
		}
	}]`
	collection = client.Database(dbName).Collection(collectionName)
	opts := options.Aggregate()
	if cur, err = collection.Aggregate(ctx, mdb.MongoPipeline(pipeline), opts); err != nil {
		fmt.Println(err)
	}
	defer cur.Close(ctx)
	count := int64(0)
	for cur.Next(ctx) {
		cur.Decode(&doc)
		fmt.Println("doc = ", doc)
		count += int64(doc["count"].(float64))
	}

	if total != count {
		fmt.Println("expected", total, "but got", count)
	}
}
