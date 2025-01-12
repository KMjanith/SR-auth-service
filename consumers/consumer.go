package consumers

import (
	services "auth-service/service"
	"auth-service/spec"
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

// Consumer1 is a simple consumer that receives a message from the queue
func Consumer1(ctx context.Context) {

	CONSUMER_CHANNEL := ctx.Value("consumerChannel").(*amqp.Channel)

	ApiAuthmsg, err := CONSUMER_CHANNEL.QueueDeclare(
		"ApiAuthMsg", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare Api-Auth-msg queue: %v", err)
	}

	msgs, err := CONSUMER_CHANNEL.Consume(
		ApiAuthmsg.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process messages
	forever := make(chan bool)

	go func() {
		for d := range msgs {

			var RegisterRequest spec.RegisterUser
			// Unmarshal the Protobuf message
			err := proto.Unmarshal(d.Body, &RegisterRequest)
			if err != nil {
				log.Fatalf("Failed to unmarshal Protobuf message: %v", err)
			}
			// Accessing fields of the protobuf message
			log.Printf("Received a RegisterRequest: username=%s", RegisterRequest.Username)
			log.Println("calling register")
			services.LoginUser(RegisterRequest.Username, RegisterRequest.Password, d.CorrelationId, d.ReplyTo, ctx)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
