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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Score struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Score int    `json:"score" bson:"score"`
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://safiqueafaruque:5AbXQyiTL05Ebnbb@cluster0.28oas.mongodb.net/scoreboard?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/scores", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		collection := client.Database("scoreboard").Collection("scores")

		cursor, err := collection.Find(context.TODO(), bson.D{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer cursor.Close(context.TODO())

		var scores []Score
		for cursor.Next(context.TODO()) {
			var score Score
			if err = cursor.Decode(&score); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = append(scores, score)
		}

		if err = cursor.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(scores)
	})

	log.Println("Server is starting...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
