package config

import (
	"os"
	"fmt"
	"log"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	srv string = os.Getenv("MSRV")
	user string = os.Getenv("MUSER")
	passw string = os.Getenv("MPASSW")
	database string = "grpc_mongodb_crud"
	collect string = "bookstore"
	client *mongo.Client

	Collection *mongo.Collection
	Mongoctx context.Context
)

func Init() {
	fmt.Println("+ Initializing database connection")
	Mongoctx = context.Background()
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf("mongodb+srv://%s:%s@%s-lvmtb.azure.mongodb.net/test?retryWrites=true&w=majority", user, passw, srv),
	)

	client, err := mongo.Connect(Mongoctx, clientOptions)
	if err != nil { log.Fatal(err) }

  	// Check whether the connection was succesful by pinging MongoDB server
	if err := client.Ping(Mongoctx, nil); err != nil {
		log.Fatalf("- Database connection error: %s", err)
	} else {
		fmt.Println("+ Database connection established")
	}

	Collection = client.Database(database).Collection(collect)
}