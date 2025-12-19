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


// Verify achievement (FR-007)
func (r *AchievementRepository) Verify(id string, verifiedBy string) error {
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
			"status": "verified",
			"verified_at": now,
			"verified_by": verifiedBy,
			"updated_at": now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Reject achievement (FR-008)
func (r *AchievementRepository) Reject(id string, reason string) error {
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
			"status": "rejected",
			"rejected_at": now,
			"rejection_reason": reason,
			"updated_at": now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}


// Find all with pagination, filtering, and sorting (FR-010)
func (r *AchievementRepository) FindAllWithPagination(filter bson.M, page, limit int, sortBy, sortOrder string) ([]model.Achievement, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Calculate skip
	skip := (page - 1) * limit

	// Determine sort direction
	sortDirection := -1 // desc by default
	if sortOrder == "asc" {
		sortDirection = 1
	}

	// Set sort options
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: sortBy, Value: sortDirection}})

	// Get total count
	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get achievements
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err = cursor.All(ctx, &achievements); err != nil {
		return nil, 0, err
	}

	return achievements, int(total), nil
}

// UpdateFields achievement
func (r *AchievementRepository) UpdateFields(id string, req *model.UpdateAchievementRequest) error {
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
			"title":       req.Title,
			"description": req.Description,
			"category":    req.Category,
			"points":      req.Points,
			"updated_at":  now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Get achievement history
func (r *AchievementRepository) GetHistory(id string) ([]model.AchievementHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Get the achievement to build history from its fields
	var achievement model.Achievement
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&achievement)
	if err != nil {
		return nil, err
	}

	// Build history from achievement data
	var history []model.AchievementHistory

	// Created event
	history = append(history, model.AchievementHistory{
		Action:    "created",
		Status:    "draft",
		Timestamp: achievement.CreatedAt,
		ActorID:   achievement.StudentID,
		ActorType: "student",
		Message:   "Achievement created",
	})

	// Submitted event
	if achievement.SubmittedAt != nil {
		history = append(history, model.AchievementHistory{
			Action:    "submitted",
			Status:    "submitted",
			Timestamp: *achievement.SubmittedAt,
			ActorID:   achievement.StudentID,
			ActorType: "student",
			Message:   "Achievement submitted for verification",
		})
	}

	// Verified event
	if achievement.VerifiedAt != nil && achievement.VerifiedBy != nil {
		history = append(history, model.AchievementHistory{
			Action:    "verified",
			Status:    "verified",
			Timestamp: *achievement.VerifiedAt,
			ActorID:   *achievement.VerifiedBy,
			ActorType: "lecturer",
			Message:   "Achievement verified by advisor",
		})
	}

	// Rejected event
	if achievement.RejectedAt != nil {
		message := "Achievement rejected"
		if achievement.RejectionReason != nil {
			message += ": " + *achievement.RejectionReason
		}
		history = append(history, model.AchievementHistory{
			Action:    "rejected",
			Status:    "rejected",
			Timestamp: *achievement.RejectedAt,
			ActorID:   achievement.StudentID, // Assuming rejection is recorded by system
			ActorType: "lecturer",
			Message:   message,
		})
	}

	return history, nil
}
// Add document to achievement
func (r *AchievementRepository) AddDocument(id string, document model.AchievementDocument) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$push": bson.M{"documents": document},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}