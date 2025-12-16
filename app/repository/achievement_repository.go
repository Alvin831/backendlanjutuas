package repository

import (
	"context"
	"time"
	"uas_backend/app/model"
	"uas_backend/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{
		collection: database.AchievementCollection,
	}
}

// Create achievement
func (r *AchievementRepository) Create(achievement *model.Achievement) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set timestamps and defaults
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	achievement.Status = "draft" // Default status
	achievement.IsDeleted = false // Default not deleted

	result, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return nil, err
	}

	// Set the ID from MongoDB
	achievement.ID = result.InsertedID.(primitive.ObjectID)
	return achievement, nil
}

// Find by ID
func (r *AchievementRepository) FindByID(id string) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var achievement model.Achievement
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

// Find by student ID
func (r *AchievementRepository) FindByStudentID(studentID string) ([]model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// Update achievement
func (r *AchievementRepository) Update(achievement *model.Achievement) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	achievement.UpdatedAt = time.Now()

	filter := bson.M{"_id": achievement.ID}
	update := bson.M{"$set": achievement}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return achievement, nil
}

// Submit achievement (change status to submitted)
func (r *AchievementRepository) Submit(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"status": "submitted",
			"submitted_at": now,
			"updated_at": now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Find all achievements (with optional filters) - exclude soft deleted
func (r *AchievementRepository) FindAll(filter bson.M) ([]model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Add filter to exclude soft deleted items
	filter["is_deleted"] = bson.M{"$ne": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// Find all achievements including deleted (for admin)
func (r *AchievementRepository) FindAllIncludeDeleted(filter bson.M) ([]model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// Find achievements by MongoDB IDs (FR-006)
func (r *AchievementRepository) FindByIDs(mongoIDs []string) ([]model.Achievement, error) {
	if len(mongoIDs) == 0 {
		return []model.Achievement{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert string IDs to ObjectIDs
	var objectIDs []primitive.ObjectID
	for _, id := range mongoIDs {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue // Skip invalid IDs
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{
		"_id": bson.M{"$in": objectIDs},
		"is_deleted": bson.M{"$ne": true},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

// Soft delete achievement (FR-005)
func (r *AchievementRepository) SoftDelete(id string, deletedBy string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": now,
			"deleted_by": deletedBy,
			"updated_at": now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Find by ID - exclude soft deleted
func (r *AchievementRepository) FindByIDActive(id string) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id": objectID,
		"is_deleted": bson.M{"$ne": true},
	}

	var achievement model.Achievement
	err = r.collection.FindOne(ctx, filter).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	return &achievement, nil
}