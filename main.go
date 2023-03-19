package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string
	Email string
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://djideath:he9SJ8TGIxKLz4qG@cluster0.y8ltv.mongodb.net/test")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Get a handle for the "users" collection
	collection := client.Database("mydb").Collection("users")

	// Set up HTTP endpoints
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(w, r, collection)
		case http.MethodPost:
			createUser(w, r, collection)
		}
	})

	http.HandleFunc("/users/{name}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// getUserByName(w, r, collection)
		case http.MethodPut:
			// updateUserByName(w, r, collection)
		case http.MethodDelete:
			// deleteUserByName(w, r, collection)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUsers(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Find all users
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	// Convert cursor to []User
	var users []User
	for cursor.Next(context.Background()) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			http.Error(w, "Failed to decode user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Convert []User to JSON and write to response
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func createUser(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	// Parse request body into User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	// Insert user into database
	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	// Return ID of inserted user
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Inserted user with ID: %v", res.InsertedID)))
}
