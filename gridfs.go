package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func GridFS() {
	var err error
	var client *mongo.Client
	var bucket *gridfs.Bucket
	var ustream *gridfs.UploadStream

	str := "This is a test file"
	client = GetMongoClient()
	defer client.Disconnect(context.Background())

	if bucket, err = gridfs.NewBucket(client.Database(dbName), options.GridFSBucket().SetName("myFiles")); err != nil {
		fmt.Println(err)
	}

	opts := options.GridFSUpload()
	opts.SetMetadata(bsonx.Doc{{Key: "content-type", Value: bsonx.String("application/json")}})
	if ustream, err = bucket.OpenUploadStream("test.txt", opts); err != nil {
		fmt.Println(err)
	}

	println("ustream = ", ustream)

	if _, err = ustream.Write([]byte(str)); err != nil {
		fmt.Println(err)
	}

	ustream.Close()
	fileID := ustream.FileID
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	if _, err = bucket.DownloadToStream(fileID, w); err != nil {
		fmt.Println(err, ustream.FileID)
	}

	fmt.Println("b.String() = ", b.String())
	if b.String() != str {
		fmt.Println("expected", str, "but got", b.String())
	}
}
