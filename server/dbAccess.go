package server

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Content string             `json:"content"`
}

var collection *mongo.Collection
var client *mongo.Client

func Connect(dbURI string) {
	clientOptions := options.Client().ApplyURI(dbURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")
	collection = client.Database("chatdemo_db").Collection("messages")
}

func Disconnect() {
	client.Disconnect(context.Background())
}

func GetMessages() ([]Message, error) {
	var messages []Message

	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var message Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func AddMessage(message Message) (Message, error) {
	insertResult, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		return message, err
	}

	message.ID = insertResult.InsertedID.(primitive.ObjectID)

	return message, nil
}
