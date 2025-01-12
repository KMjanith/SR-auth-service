package producers

import (
	"auth-service/spec"
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func SendRegisterResponse(response string, token string, userID string, corrID string,replyTo string,  ctx context.Context) {
	msg := &spec.RegisterUserResponse{
		UserId: userID,
		Message: response,
		Token: token,
	}

	ch := ctx.Value("producerChannel")
	

	// Serialize message using protobuf
	request, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.(*amqp.Channel).PublishWithContext(ctx,
		"",     // exchange
		replyTo, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/protobuf",
			CorrelationId: corrID,
			Body:        request,
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
	log.Printf(" [x] Sent %s\n", msg)
}