package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
}

func New() *Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	c := &Client{
		client: client,
	}

	users := c.users()
	if _, err := users.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		panic(err)
	}

	return c
}

func (c *Client) users() *mongo.Collection {
	return c.collection("users")
}

func (c *Client) collection(collection string) *mongo.Collection {
	return c.db().Collection(collection)
}

func (c *Client) db() *mongo.Database {
	return c.client.Database("db")
}
