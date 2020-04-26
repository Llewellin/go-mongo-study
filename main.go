package main

import (
	"context"
	"io"

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

// var session *mgo.Session

type currency struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Amount  float64            `bson:"amount"`
	Account string             `bson:"account"`
	Code    string             `bson:"code"`
	Version int                `bson:"version"`
}

// Random get random value
func Random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min+1) + min
}

func pay(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("???")
	// session.StartTransaction
	// entry := currency{}
	var findResult *mongo.SingleResult
	var err error
	var entry currency
	// fmt.Println("globalDB = ", globalDB)

	//step 3: subtract current balance and update back to database
	// err = globalDB.C("bank").UpdateId(entry.ID, &entry)
	// var result *mongo.UpdateResult
	// var result2 mongo.UpdateResult

	// fmt.Println("result = ", result)
	// fmt.Println("result2 = ", result2)

	// step 1: get current amount
	findResult = globalDB.Collection("bank").FindOne(ctx, bson.M{"account": account})

	err = findResult.Err()
	if err != nil {
		// panic("1")
		fmt.Println("err = ", err)
		return
		// panic(err)
	}

	err = findResult.Decode(&entry)
	if err != nil {
		// panic("2")
		fmt.Println("2")
		return
		// panic(err)
	}
	wait := Random(1, 100)
	time.Sleep(time.Duration(wait) * time.Millisecond)

	entry.Amount = entry.Amount + 1.000

	// fmt.Println("version = ", entry.Version)

	_, err = globalDB.Collection("bank").UpdateOne(ctx, bson.M{"_id": entry.ID}, bson.M{"$set": bson.M{
		"amount": entry.Amount,
	}})
	// fmt.Println("UpdateOne() result MatchedCount:", result.MatchedCount)
	// fmt.Println("UpdateOne() result ModifiedCount:", result.ModifiedCount)

	// }

	// var wtf interface{}
	// err = updateResult.Decode(&wtf)
	// fmt.Println("UpdateOne() result:", *result)
	// fmt.Println("UpdateOne() result TYPE:", reflect.TypeOf(result))
	// fmt.Println("UpdateOne() result MatchedCount:", result.MatchedCount)
	// fmt.Println("UpdateOne() result UpsertedCount:", result.UpsertedCount)
	// fmt.Println("UpdateOne() result UpsertedID:", result.UpsertedID)

	if err != nil {
		// panic("update error")
		fmt.Println("3")
		return
	}

	fmt.Printf("%+v\n", entry)

	// haha := globalDB.Collection("bank").FindOneAndUpdate(ctx, bson.M{"account": account}, bson.M{"$set": bson.M{
	// 	"amount": entry.Amount + 1,
	// }})
	// var entry2 currency
	// err = haha.Decode(&entry2)
	// if err != nil {
	// 	// panic("2")
	// 	fmt.Println("2")
	// 	return
	// 	// panic(err)
	// }

	// fmt.Println("entry2 =", entry2)

	io.WriteString(w, "ok")
}

func main() {
	var err error
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	// Base context.
	ctx = context.Background()
	clientOpts := options.Client().ApplyURI("mongodb://@localhost:27017/")
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

	coll.DeleteMany(ctx, bson.M{})

	fmt.Println("globalDB = ", globalDB)

	var result *mongo.InsertOneResult
	user := currency{Account: account, Amount: 1000.00, Code: "USD", Version: 0}
	result, err = coll.InsertOne(ctx, user)

	if err != nil {
		panic("insert error")
	}

	fmt.Println("result = ", result)

	log.Println("Listen server on " + port + " port")
	http.HandleFunc("/", pay)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
