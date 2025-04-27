package StorageResources

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client      *mongo.Client
	Departments *mongo.Collection
	Curriculums *mongo.Collection
	Subjects    *mongo.Collection
	Rooms       *mongo.Collection
	Instructors *mongo.Collection
}

func NewMongodbClient() *mongo.Client {

	//////////////////////////////////////////////////////////////////////////
	// MongoDB Setup
	//////////////////////////////////////////////////////////////////////////

	log.Println("connecting to MongoDB...")

	log.Printf("MONGO_DB_USER     = %s\n", os.Getenv("MONGO_DB_USER"))
	log.Printf("MONGO_DB_PASSWORD = %s\n", os.Getenv("MONGO_DB_PASSWORD"))
	log.Printf("PORT              = %s\n", os.Getenv("PORT"))

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb+srv://%s:%s@testcluster.sz6qg.mongodb.net/?retryWrites=true&w=majority&appName=TestCluster",
		os.Getenv("MONGO_DB_USER"),
		os.Getenv("MONGO_DB_PASSWORD"),
	)).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	log.Println("connected to MongoDB...")

	return client
}

func CloseMongodbClient(client *mongo.Client) error {
	err := client.Disconnect(context.TODO())
	log.Print(err)
	return err
}
