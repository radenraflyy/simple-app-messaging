package repository

import (
	"SimpleMessaging/app/models"
	"SimpleMessaging/pkg/database"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertNewMessage(ctx context.Context, data models.MessagePayload) error {
	_, err := database.MongoDB.InsertOne(ctx, data)
	return fmt.Errorf("failed to insert new message: %v", err)
}

func GetAllMessage(ctx context.Context) ([]models.MessagePayload, error) {
	var (
		err  error
		resp []models.MessagePayload
	)

	cursor, err := database.MongoDB.Find(ctx, bson.D{})
	if err != nil {
		return resp, fmt.Errorf("failed to find message: %v", err)
	}
	for cursor.Next(ctx) {
		payload := models.MessagePayload{}
		if err := cursor.Decode(&payload); err != nil {
			return resp, fmt.Errorf("failed to decode message: %v", err)
		}
		resp = append(resp, payload)
	}
	return resp, nil
}
