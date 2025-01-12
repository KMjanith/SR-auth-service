package main

import (
	"auth-service/consumers"
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	godotenv.Load() //make the connection to the mongodb
	log.Println("starting Auth service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Set up a context with timeout
	defer cancel()

	//----------------------------------------MONGO CONNECTION----------------------------------//

	serverAPI := options.ServerAPI(options.ServerAPIVersion1) // Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	mongoURI := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts) // Create a new client and connect to the server
	if err != nil {
		log.Printf("Error:%s ", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Printf("Error: %s", err)
		}
	}()
	var result bson.M // Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
	ctx = context.WithValue(ctx, "mongoClient", client) // Store client in context

	//---------------------------------------RABBITMQ CONNECTION---------------------------------//

	rabbitUrl := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(rabbitUrl) //making the connection to the rabbitmq
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	ctx = context.WithValue(ctx, "conn", conn) //adding the connection to the context
	consumerChannel, err := conn.Channel()     //Open a channel for consumer
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer consumerChannel.Close()
	ctx = context.WithValue(ctx, "consumerChannel", consumerChannel) //adding the consumer channel to the context
	producerChannel, err := conn.Channel()                           //open a channel for producer
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer producerChannel.Close()
	log.Println("Connected to RabbitMQ")
	ctx = context.WithValue(ctx, "producerChannel", producerChannel) //adding the producer channel to the context

	consumers.Consumer1(ctx) //run the consumer

}
