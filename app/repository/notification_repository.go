package repository

import (
	"context"
	"time"
	"uas_backend/app/model"
	"uas_backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository() *NotificationRepository {
	// Get notifications collection from MongoDB
	notificationCollection := database.MongoClient.Database("achievement_db").Collection("notifications")
	
	return &NotificationRepository{
		collection: notificationCollection,
	}
}

// Create notification
func (r *NotificationRepository) Create(notification *model.Notification) (*model.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	notification.CreatedAt = time.Now()
	notification.IsRead = false

	result, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return nil, err
	}

	notification.ID = result.InsertedID.(primitive.ObjectID)
	return notification, nil
}

// Find notifications by recipient ID
func (r *NotificationRepository) FindByRecipientID(recipientID string, limit int) ([]model.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"recipient_id": recipientID}
	
	// Sort by created_at descending, limit results
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []model.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Mark notification as read
func (r *NotificationRepository) MarkAsRead(notificationID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return err
	}

	now := time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
			"read_at": now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}