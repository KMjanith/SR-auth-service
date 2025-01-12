package services

import (
	"auth-service/producers"
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Struct for user registration data
type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register a new user
func LoginUser(userName string, password string, corrID string,replyTo string, ctx context.Context) {

	var newUser User
	newUser.Username = userName
	newUser.Password = password

	MONGO_CLIENT := ctx.Value("mongoClient").(*mongo.Client)
	coll := MONGO_CLIENT.Database(os.Getenv("DATABASE")).Collection(os.Getenv("USER_COLLECTION"))

	// Finding if the user already exists
	var result User
	filter := bson.D{{Key: "username", Value: newUser.Username}}
	err := coll.FindOne(context.Background(), filter).Decode(&result)

	log.Printf("username: %s, password: %s", userName, password)

	// If user is found, check the password
	if err == nil {
		log.Println("user already exist")
		// User found, check if the password is correct
		err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(newUser.Password))
		if err == nil {
			log.Println(err)
			// Password is correct, generate token and send success response
			tokenString, err := GenerateToken(newUser.Username, result.ID)
			log.Println(tokenString)
			if err != nil {
				log.Println("Error generating token:", err)
				producers.SendRegisterResponse("Internal server error", "", "", corrID, replyTo,ctx)
				return
			}
			producers.SendRegisterResponse("Password is correct", tokenString, result.ID, corrID, replyTo,ctx)
			return
		} else {
			// Password is incorrect, update credentials
			log.Println("Incorrect password, updating credentials")

			// Delete the existing user
			_, err := coll.DeleteOne(context.Background(), bson.D{{Key: "username", Value: newUser.Username}})
			if err != nil {
				log.Fatal("Error deleting user:", err)
				producers.SendRegisterResponse(err.Error(), "", "", corrID, replyTo,ctx)
				return
			}
		}
	}

	// Hash the password for new or updated user
	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		log.Fatal(err)
		producers.SendRegisterResponse(err.Error(), "", "", corrID, replyTo,ctx)
	}

	newUser.Password = hashedPassword

	// Insert the newUser into the database.
	insertResult, err := coll.InsertOne(context.Background(), newUser)
	if err != nil {
		log.Fatal(err)
		producers.SendRegisterResponse(err.Error(), "", "", corrID,replyTo, ctx)
	}

	// Retrieve the newly inserted document ID.
	insertedID := insertResult.InsertedID
	log.Printf("Inserted a single document with ID: %v\n", insertedID)

	// Convert the insertedID to a string if it's an ObjectID.
	idStr := fmt.Sprintf("%v", insertedID)

	// Token generation
	tokenString, err := GenerateToken(newUser.Username, idStr)
	if err != nil {
		producers.SendRegisterResponse("Internal server error", "", "", corrID,replyTo, ctx)
		return
	}

	// Respond with the new user and the token.
	producers.SendRegisterResponse("User registered successfully", tokenString, idStr, corrID,replyTo, ctx)
}

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
