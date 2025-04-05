package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"login-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUpHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		collection := client.Database("your_database").Collection("users")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if the email already exists
		var existingUser models.User
		err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}

		// Insert user into MongoDB
		_, err = collection.InsertOne(ctx, user)
		if err != nil {
			log.Println("Error inserting user:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	}
}
