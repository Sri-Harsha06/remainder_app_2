package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"remainder_app_2/model"
	"remainder_app_2/services"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func Findtmrevents(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var events []model.Event
	client = services.GetInstance()
	collection := client.Database("events").Collection("event")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	dt := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	cursor, err := collection.Find(ctx, model.Event{Date: dt})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	collection = client.Database("tmrevents").Collection("event")
	_, err = collection.DeleteMany(context.Background(), model.Event{})
	for cursor.Next(ctx) {
		var event model.Event
		cursor.Decode(&event)
		events = append(events, event)
		_, _ = collection.InsertOne(context.Background(), event)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(events)
}
