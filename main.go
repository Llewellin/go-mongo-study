package main

import (
	"context"

	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	// "github.com/globalsign/mgo"
	// "github.com/globalsign/mgo/bson"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var globalDB *mongo.Database
var account = "vincent"
var ctx context.Context
var client *mongo.Client
var session mongo.Session

var dbName = "user"
var collectionExamples = "account"

// var session *mgo.Session

type currency struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount  float64            `bson:"amount"`
	Account string             `bson:"account"`
	Code    string             `bson:"code"`
	Version int                `bson:"version"`
}

func TestTransactionCommit() {
	var err error
	// var client = client
	var collection *mongo.Collection
	var ctx = context.Background()
	var id = primitive.NewObjectID()
	var doc = bson.M{"_id": id, "hometown": "Atlanta", "year": int32(1998)}
	var result *mongo.UpdateResult
	var session mongo.Session
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: int32(2000)}}}}
	// client = client
	defer client.Disconnect(ctx)
	collection = client.Database(dbName).Collection(collectionExamples)
	if _, err = collection.InsertOne(ctx, doc); err != nil {
		fmt.Println(err)
	}

	if session, err = client.StartSession(); err != nil {
		fmt.Println(err)
	}
	if err = session.StartTransaction(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("start transaction")
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if result, err = collection.UpdateOne(sc, bson.M{"_id": id}, update); err != nil {
			fmt.Println(err)
		}
		if result.MatchedCount != 1 || result.ModifiedCount != 1 {
			fmt.Println("replace failed, expected 1 but got", result.MatchedCount)
		}

		if err = session.CommitTransaction(sc); err != nil {
			fmt.Println(err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	session.EndSession(ctx)

	var v bson.M
	if err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
		fmt.Println(err)
	}
	fmt.Println("v = ", v)
	// if v["year"] != int32(2000) {
	// 	// fmt.Println(stringify(v))
	// 	fmt.Println("expected 2000 but got", v["year"])
	// }

	// res, _ := collection.DeleteOne(ctx, bson.M{"_id": id})
	// if res.DeletedCount != 1 {
	// 	fmt.Println("delete failed, expected 1 but got", res.DeletedCount)
	// }
}

func abort(w http.ResponseWriter, r *http.Request) {
	TestTransactionAbort()
}

func TestTransactionAbort() {
	var err error
	// var client *mongo.Client
	var collection *mongo.Collection
	var ctx = context.Background()
	var id = primitive.NewObjectID()
	var id2 = primitive.NewObjectID()
	// var doc = currency{Account: account, Amount: 1000.00, Code: "USD", Version: 0}
	var doc = bson.M{"_id": id, "hometown": "Atlanta", "year": int32(1998)}
	var doc2 = bson.M{"_id": id2, "hometown": "Universe", "year": int32(1990)}
	var result *mongo.UpdateResult
	var session mongo.Session
	var update = bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: int32(2000)}}}}
	// client = client
	defer client.Disconnect(ctx)
	collection = client.Database(dbName).Collection(collectionExamples)

	if session, err = client.StartSession(); err != nil {
		fmt.Println(err)
	}
	if err = session.StartTransaction(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("ctx = ", ctx)
	fmt.Println("session = ", session)
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		fmt.Println("sc = ", sc)

		if _, err = collection.InsertOne(sc, doc); err != nil {
			fmt.Println(err)
		}
		if _, err = collection.InsertOne(sc, doc2); err != nil {
			fmt.Println(err)
		}
		if result, err = collection.UpdateOne(sc, bson.M{"_id": id}, update); err != nil {
			fmt.Println(err)
		}
		if _, err = collection.InsertOne(sc, doc); err != nil {
			fmt.Println("error occur")
			if err = session.AbortTransaction(sc); err != nil {
				fmt.Println(err)
			}
			fmt.Println(err)
		}
		if result.MatchedCount != 1 || result.ModifiedCount != 1 {
			fmt.Println("replace failed, expected 1 but got", result.MatchedCount)
		}

		if err = session.AbortTransaction(sc); err != nil {
			fmt.Println(err)
		}
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	session.EndSession(ctx)

	// wait := Random(100, 1000)
	// time.Sleep(time.Duration(wait) * time.Millisecond)
	var v bson.M
	if err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&v); err != nil {
		fmt.Println(err)
	}
	fmt.Println("v = ", v)
	// if v["year"] != int32(1998) {
	// 	// t.Log(stringify(v))
	// 	fmt.Println("expected 1998 but got", v["year"])
	// }

	// res, _ := collection.DeleteOne(ctx, bson.M{"_id": id})
	// if res.DeletedCount != 1 {
	// 	fmt.Println("delete failed, expected 1 but got", res.DeletedCount)
	// }
}

// Random get random value
func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}

func pay(w http.ResponseWriter, r *http.Request) {
	// // fmt.Println("???")
	// // session.StartTransaction
	// // entry := currency{}
	// var findResult *mongo.SingleResult
	// var err error
	// var entry currency
	// // fmt.Println("globalDB = ", globalDB)

	// //step 3: subtract current balance and update back to database
	// // err = globalDB.C("bank").UpdateId(entry.ID, &entry)
	// // var result *mongo.UpdateResult
	// // var result2 mongo.UpdateResult

	// // fmt.Println("result = ", result)
	// // fmt.Println("result2 = ", result2)

	// // step 1: get current amount
	// findResult = globalDB.Collection("bank").FindOne(ctx, bson.M{"account": account})

	// err = findResult.Err()
	// if err != nil {
	// 	// panic("1")
	// 	fmt.Println("err = ", err)
	// 	return
	// 	// panic(err)
	// }

	// err = findResult.Decode(&entry)
	// if err != nil {
	// 	// panic("2")
	// 	fmt.Println("2")
	// 	return
	// 	// panic(err)
	// }
	// wait := Random(1, 100)
	// time.Sleep(time.Duration(wait) * time.Millisecond)

	// entry.Amount = entry.Amount + 1.000

	// // fmt.Println("version = ", entry.Version)

	// _, err = globalDB.Collection("bank").UpdateOne(ctx, bson.M{"_id": entry.ID}, bson.M{"$set": bson.M{
	// 	"amount": entry.Amount,
	// }})
	// // fmt.Println("UpdateOne() result MatchedCount:", result.MatchedCount)
	// // fmt.Println("UpdateOne() result ModifiedCount:", result.ModifiedCount)

	// // }

	// // var wtf interface{}
	// // err = updateResult.Decode(&wtf)
	// // fmt.Println("UpdateOne() result:", *result)
	// // fmt.Println("UpdateOne() result TYPE:", reflect.TypeOf(result))
	// // fmt.Println("UpdateOne() result MatchedCount:", result.MatchedCount)
	// // fmt.Println("UpdateOne() result UpsertedCount:", result.UpsertedCount)
	// // fmt.Println("UpdateOne() result UpsertedID:", result.UpsertedID)

	// if err != nil {
	// 	// panic("update error")
	// 	fmt.Println("3")
	// 	return
	// }

	// fmt.Printf("%+v\n", entry)

	// // haha := globalDB.Collection("bank").FindOneAndUpdate(ctx, bson.M{"account": account}, bson.M{"$set": bson.M{
	// // 	"amount": entry.Amount + 1,
	// // }})
	// // var entry2 currency
	// // err = haha.Decode(&entry2)
	// // if err != nil {
	// // 	// panic("2")
	// // 	fmt.Println("2")
	// // 	return
	// // 	// panic(err)
	// // }

	// // fmt.Println("entry2 =", entry2)

	// io.WriteString(w, "ok")

	TestTransactionCommit()
}

func main() {
	var err error
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// Base context.
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	clientOpts := options.Client().ApplyURI("mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=rs0")
	client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		fmt.Println(err)
		return
	}

	globalDB = client.Database("queue")
	coll := globalDB.Collection("bank")
	// session, _ = mgo.Dial("localhost:27017")
	// globalDB = session.DB("queue")
	// globalDB.C("bank").DropCollection()
	fmt.Println("globalDB = ", globalDB)

	deleteResult, err := coll.DeleteMany(ctx, bson.M{})
	fmt.Println("deleteResult = ", deleteResult)
	if err != nil {
		fmt.Println("err delete = ", err)

	}

	var result *mongo.InsertOneResult
	user := currency{Account: account, Amount: 1000.00, Code: "USD", Version: 0}
	result, err = coll.InsertOne(ctx, user)

	if err != nil {
		fmt.Println("err = ", err)
		panic("insert error")
	}

	fmt.Println("result = ", result)

	log.Println("Listen server on " + port + " port")
	http.HandleFunc("/pay", pay)
	http.HandleFunc("/abort", abort)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
