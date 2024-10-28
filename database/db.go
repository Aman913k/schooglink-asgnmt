package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://amanrana9133:aman123@cluster0.6ohqxgo.mongodb.net/?retryWrites=true&w=majority"
const dbName = "schooglink"


var collection *mongo.Collection

func GetCollection(colName string) *mongo.Collection {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo Connection Successfull")
	collection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")

	return collection

}
