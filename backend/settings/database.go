package settings

import (
	"context"
	"net/url"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database
var PageSize string = "10"

func Create_database_client() {
	MongoClient = establish_mongodb_connection()
	MongoDatabase = access_mongodb_database(MongoClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Cannot connect to database server, please check the configuration or status of the database server")
	}
	log.Info("Completed setting up a database client")

}

func establish_mongodb_connection() *mongo.Client {
	// Establish database connection
	_, existed := Evariables["MONGODB_URI"]
	if !existed {
		username := Evariables["MONGODB_USER"]
		password := Evariables["MONGODB_PASSWORD"]
		cluster := Evariables["MONGODB_SERVER_URL"]
		uri := "mongodb+srv://" + url.QueryEscape(username) + ":" +
			url.QueryEscape(password) + "@" + cluster
		log.Info(uri)
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

		if err != nil {
			log.Fatal("Failed to connect to the database server, please check the URL variable")
		}
		return client
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(Evariables["MONGODB_URI"]))
	if err != nil {
		log.Fatal("Failed to connect to the database server, please check the URL variable")
	}

	return client
}

func access_mongodb_database(client *mongo.Client) *mongo.Database {
	_, existed := Evariables["MONGODB_DATABASE"]
	if !existed {
		log.Fatal("The variable MONGODB_DATABASE doesn't exist")
	}
	database := client.Database(Evariables["MONGODB_DATABASE"])
	return database
}
