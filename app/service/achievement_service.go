package service

import (
	"context"
	"fmt"
	"time"
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"uas_backend/app/utils"
	"uas_backend/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var achievementRepo *repository.AchievementRepository
var achievementRefRepo *repository.AchievementReferenceRepository
var studentRepo *repository.StudentRepository

func SetAchievementRepo(repo *repository.AchievementRepository) {
	achievementRepo = repo
}

func SetAchievementReferenceRepo(repo *repository.AchievementReferenceRepository) {
	achievementRefRepo = repo
}

func SetStudentRepo(repo *repository.StudentRepository) {
	studentRepo = repo
}

// GetAllAchievements godoc
// @Summary Get All Achievements
// @Description Get list of achievements with role-based filtering and pagination
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status" Enums(draft, submitted, verified, rejected)
// @Param category query string false "Filter by category"
// @Param student_id query string false "Filter by student ID"
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} utils.Response "Achievements retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/achievements [get]
func GetAllAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	// Get pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	
	// Get filter parameters
	status := c.Query("status")        // draft, submitted, verified, rejected
	category := c.Query("category")    // category filter
	studentID := c.Query("student_id") // filter by student
	
	// Get sort parameters
	sortBy := c.Query("sort_by", "created_at") // created_at, updated_at, points
	sortOrder := c.Query("sort_order", "desc") // asc, desc
	
	var filter bson.M
	
	// Role-based filtering
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" { // Mahasiswa role ID
		// Mahasiswa hanya bisa lihat achievement sendiri
		filter = bson.M{"student_id": userID, "is_deleted": bson.M{"$ne": true}}
	} else {
		// Admin/Dosen bisa lihat semua
		filter = bson.M{"is_deleted": bson.M{"$ne": true}}
		
		// Apply filters (only for admin/dosen)
		if status != "" {
			filter["status"] = status
		}
		if category != "" {
			filter["category"] = category
		}
		if studentID != "" {
			filter["student_id"] = studentID
		}
	}
	
	// Get achievements with pagination and sorting
	achievements, total, err := achievementRepo.FindAllWithPagination(filter, page, limit, sortBy, sortOrder)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get achievements", 500, nil))
	}
	
	// Calculate pagination info
	totalPages := (total + limit - 1) / limit
	
	return c.JSON(utils.SuccessResponse("Achievements retrieved", 200, fiber.Map{
		"achievements": achievements,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
		"filters": fiber.Map{
			"status":     status,
			"category":   category,
			"student_id": studentID,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
		},
	}))
}

// GET /api/v1/achievements/:id - Detail
func GetAchievementByID(c *fiber.Ctx) error {
	id := c.Params("id")
	
	// Use active finder to exclude soft deleted achievements
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// Check access: mahasiswa hanya bisa lihat achievement sendiri
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" && achievement.StudentID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Access denied", 403, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Achievement retrieved", 200, achievement))
}

// CreateAchievement godoc
// @Summary Create New Achievement
// @Description Create a new achievement (Mahasiswa only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateAchievementRequest true "Achievement data"
// @Success 201 {object} utils.Response "Achievement created successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/achievements [post]
func CreateAchievement(c *fiber.Ctx) error {
	var req model.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// Get user info from JWT
	userID := c.Locals("user_id").(string)
	
	// Create achievement object
	achievement := &model.Achievement{
		StudentID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Points:      req.Points,
		Status:      "draft", // FR-003: Status awal 'draft'
		Documents:   []model.AchievementDocument{}, // Empty initially
	}
	
	// Save to MongoDB
	createdAchievement, err := achievementRepo.Create(achievement)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to create achievement", 500, nil))
	}
	
	// Create reference in PostgreSQL
	if achievementRefRepo != nil && studentRepo != nil {
		// Get student ID from user ID
		student, err := studentRepo.FindByUserID(userID)
		if err == nil && student != nil {
			ref := &model.AchievementReference{
				MongoID:   createdAchievement.ID.Hex(),
				StudentID: student.ID, // Use student.id instead of user_id
				Status:    "draft",
				CreatedAt: createdAchievement.CreatedAt,
				UpdatedAt: createdAchievement.UpdatedAt,
			}
			
			err = achievementRefRepo.Create(ref)
			if err != nil {
				// Log error but continue
				fmt.Printf("Failed to create PostgreSQL reference: %v\n", err)
			}
		}
	}
	
	// FR-003: Return achievement data
	return c.Status(201).JSON(utils.SuccessResponse("Achievement created successfully", 201, createdAchievement))
}

// PUT /api/v1/achievements/:id - Update (Mahasiswa)
func UpdateAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// Get achievement by ID first
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// Check ownership - mahasiswa hanya bisa update achievement sendiri, admin bisa update semua
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	// Mahasiswa hanya bisa update achievement sendiri
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" && achievement.StudentID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Access denied - not your achievement", 403, nil))
	}
	
	// Check precondition - hanya achievement dengan status 'draft' yang bisa diupdate
	if achievement.Status != "draft" {
		return c.Status(400).JSON(utils.ErrorResponse("Only draft achievements can be updated", 400, nil))
	}
	
	// Update achievement
	err = achievementRepo.UpdateFields(id, &req)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update achievement", 500, nil))
	}
	
	// Get updated achievement
	updatedAchievement, err := achievementRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get updated achievement", 500, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Achievement updated successfully", 200, updatedAchievement))
}

// DELETE /api/v1/achievements/:id - Delete (Mahasiswa) - FR-005 Implementation
func DeleteAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	
	// 1. Get achievement by ID (using active finder to exclude already deleted)
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// 2. Check ownership - mahasiswa hanya bisa delete achievement sendiri, admin bisa delete semua
	role := c.Locals("role").(string)
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" && achievement.StudentID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Access denied - not your achievement", 403, nil))
	}
	
	// 3. Check precondition - prestasi harus berstatus 'draft'
	if achievement.Status != "draft" {
		return c.Status(400).JSON(utils.ErrorResponse("Only draft achievements can be deleted", 400, nil))
	}
	
	// 4. Soft delete data di MongoDB
	err = achievementRepo.SoftDelete(id, userID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to delete achievement", 500, nil))
	}
	
	// 5. Update reference di PostgreSQL (if exists)
	if achievementRefRepo != nil {
		err = achievementRefRepo.SoftDelete(id)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to update PostgreSQL reference: %v\n", err)
		}
	}
	
	// 6. Return success message
	return c.JSON(utils.SuccessResponse("Achievement deleted successfully", 200, fiber.Map{
		"id": id,
		"message": "Achievement has been deleted and moved to trash",
	}))
}

// SubmitAchievement godoc
// @Summary Submit Achievement for Verification
// @Description Submit a draft achievement for verification by advisor
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param request body model.SubmitAchievementRequest true "Submit data"
// @Success 200 {object} utils.Response "Achievement submitted successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 403 {object} utils.Response "Access denied"
// @Failure 404 {object} utils.Response "Achievement not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/achievements/{id}/submit [post]
func SubmitAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	
	var req model.SubmitAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// 1. Get achievement by ID (using active finder)
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// 2. Check ownership - mahasiswa hanya bisa submit achievement sendiri, admin bisa submit semua
	role := c.Locals("role").(string)
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" && achievement.StudentID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Access denied - not your achievement", 403, nil))
	}
	
	// 3. Check precondition - prestasi harus berstatus 'draft'
	if achievement.Status != "draft" {
		return c.Status(400).JSON(utils.ErrorResponse("Achievement must be in 'draft' status to submit", 400, nil))
	}
	
	// 4. Update status menjadi 'submitted'
	err = achievementRepo.Submit(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to submit achievement", 500, nil))
	}
	
	// Update reference in PostgreSQL
	if achievementRefRepo != nil {
		ref := &model.AchievementReference{
			MongoID:   id,
			Status:    "submitted",
			UpdatedAt: time.Now(),
		}
		
		err = achievementRefRepo.Update(ref)
		if err != nil {
			fmt.Printf("Failed to update PostgreSQL reference: %v\n", err)
		}
	}
	
	// 5. Get updated achievement
	updatedAchievement, err := achievementRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get updated achievement", 500, nil))
	}
	
	// 6. Create notification untuk dosen wali
	// TODO: Get advisor ID from student data (for now using mock advisor ID)
	advisorID := "mock-advisor-id" // In real implementation, get from student-advisor relationship
	
	err = CreateAchievementSubmissionNotification(
		advisorID,
		userID,
		id,
		achievement.Title,
	)
	if err != nil {
		// Log error but don't fail the submission
		fmt.Printf("Failed to create notification: %v\n", err)
	}
	
	// 7. Return updated status
	return c.JSON(utils.SuccessResponse("Achievement submitted for verification successfully", 200, fiber.Map{
		"achievement": updatedAchievement,
		"message": "Your achievement has been submitted and is now waiting for verification by your advisor",
	}))
}

// VerifyAchievement godoc
// @Summary Verify Achievement
// @Description Verify a submitted achievement (Dosen Wali only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Achievement ID"
// @Param request body model.VerifyAchievementRequest true "Verification data"
// @Success 200 {object} utils.Response "Achievement verified successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 404 {object} utils.Response "Achievement not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/achievements/{id}/verify [post]
func VerifyAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	
	var req model.VerifyAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// 1. Get achievement by ID (using active finder)
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// 2. Check precondition - prestasi harus berstatus 'submitted'
	if achievement.Status != "submitted" {
		return c.Status(400).JSON(utils.ErrorResponse("Achievement must be in 'submitted' status to verify", 400, nil))
	}
	
	// 3. Update status menjadi 'verified' dengan verified_by dan verified_at
	err = achievementRepo.Verify(id, userID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to verify achievement", 500, nil))
	}
	
	// 4. Update reference in PostgreSQL
	if achievementRefRepo != nil {
		ref := &model.AchievementReference{
			MongoID:   id,
			Status:    "verified",
			UpdatedAt: time.Now(),
		}
		
		err = achievementRefRepo.Update(ref)
		if err != nil {
			fmt.Printf("Failed to update PostgreSQL reference: %v\n", err)
		}
	}
	
	// 5. Get updated achievement
	updatedAchievement, err := achievementRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get updated achievement", 500, nil))
	}
	
	// 6. Create notification untuk mahasiswa
	err = CreateAchievementVerifiedNotification(
		achievement.StudentID,
		userID,
		id,
		achievement.Title,
	)
	if err != nil {
		// Log error but don't fail the verification
		fmt.Printf("Failed to create notification: %v\n", err)
	}
	
	// 7. Return updated status
	return c.JSON(utils.SuccessResponse("Achievement verified successfully", 200, fiber.Map{
		"achievement": updatedAchievement,
		"message": "Achievement has been verified",
	}))
}

// POST /api/v1/achievements/:id/reject - Reject (Dosen Wali) - FR-008 Implementation
func RejectAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(string)
	
	var req model.RejectAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// Validate reason is provided
	if req.Reason == "" {
		return c.Status(400).JSON(utils.ErrorResponse("Rejection reason is required", 400, nil))
	}
	
	// 1. Get achievement by ID (using active finder)
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// 2. Check precondition - prestasi harus berstatus 'submitted'
	if achievement.Status != "submitted" {
		return c.Status(400).JSON(utils.ErrorResponse("Achievement must be in 'submitted' status to reject", 400, nil))
	}
	
	// 3. Update status menjadi 'rejected' dengan rejection reason
	err = achievementRepo.Reject(id, req.Reason)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to reject achievement", 500, nil))
	}
	
	// 4. Update reference in PostgreSQL
	if achievementRefRepo != nil {
		ref := &model.AchievementReference{
			MongoID:   id,
			Status:    "rejected",
			UpdatedAt: time.Now(),
		}
		
		err = achievementRefRepo.Update(ref)
		if err != nil {
			fmt.Printf("Failed to update PostgreSQL reference: %v\n", err)
		}
	}
	
	// 5. Get updated achievement
	updatedAchievement, err := achievementRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get updated achievement", 500, nil))
	}
	
	// 6. Create notification untuk mahasiswa
	err = CreateAchievementRejectedNotification(
		achievement.StudentID,
		userID,
		id,
		achievement.Title,
		req.Reason,
	)
	if err != nil {
		// Log error but don't fail the rejection
		fmt.Printf("Failed to create notification: %v\n", err)
	}
	
	// 7. Return updated status
	return c.JSON(utils.SuccessResponse("Achievement rejected", 200, fiber.Map{
		"achievement": updatedAchievement,
		"message": "Achievement has been rejected",
		"reason": req.Reason,
	}))
}

// GET /api/v1/achievements/:id/history - Status history
func GetAchievementHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	// Get achievement by ID first to check access
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// Check access: mahasiswa hanya bisa lihat history achievement sendiri
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" && achievement.StudentID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Access denied", 403, nil))
	}
	
	// Get achievement history
	history, err := achievementRepo.GetHistory(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get achievement history", 500, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Achievement history retrieved", 200, fiber.Map{
		"achievement_id": id,
		"history": history,
	}))
}

// POST /api/v1/achievements/:id/attachments - Upload files
func UploadAchievementAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	// Get achievement by ID first to check access
	achievement, err := achievementRepo.FindByIDActive(id)
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse("Achievement not found", 404, nil))
	}
	
	// Check ownership - mahasiswa hanya bisa upload ke achievement sendiri, admin bisa upload ke semua
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" && achievement.StudentID != userID {
		return c.Status(403).JSON(utils.ErrorResponse("Access denied - not your achievement", 403, nil))
	}
	
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("No file uploaded", 400, nil))
	}
	
	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return c.Status(400).JSON(utils.ErrorResponse("File size too large (max 10MB)", 400, nil))
	}
	
	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/jpg":  true,
		"application/pdf": true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(400).JSON(utils.ErrorResponse("File type not allowed", 400, nil))
	}
	
	// Create uploads directory if not exists
	uploadsDir := "uploads/achievements"
	if err := c.SaveFile(file, uploadsDir+"/"+file.Filename); err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to save file", 500, nil))
	}
	
	// Create document record
	document := model.AchievementDocument{
		FileName:    file.Filename,
		FileSize:    file.Size,
		ContentType: contentType,
		FilePath:    uploadsDir + "/" + file.Filename,
		UploadedAt:  time.Now(),
	}
	
	// Add document to achievement
	err = achievementRepo.AddDocument(id, document)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to add document to achievement", 500, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Attachment uploaded successfully", 200, fiber.Map{
		"achievement_id": id,
		"document": document,
	}))
}

// GET /api/v1/achievements/advisees - View Prestasi Mahasiswa Bimbingan (FR-006)
func GetAdviseesAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	// Get pagination parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit
	
	// Check if repositories are initialized
	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}
	
	if achievementRefRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Achievement reference repository not initialized", 500, nil))
	}
	
	if achievementRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Achievement repository not initialized", 500, nil))
	}
	
	// 1. Get list student IDs dari tabel students where advisor_id
	studentIDs, err := studentRepo.GetStudentIDsByAdvisorID(userID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get advisees: " + err.Error(), 500, nil))
	}
	
	if len(studentIDs) == 0 {
		return c.JSON(utils.SuccessResponse("No advisees found", 200, fiber.Map{
			"achievements": []interface{}{},
			"students": []interface{}{},
			"pagination": fiber.Map{
				"page": page,
				"limit": limit,
				"total": 0,
				"total_pages": 0,
			},
		}))
	}
	
	// 2. Get achievements references dengan filter student_ids
	achievementRefs, err := achievementRefRepo.GetByStudentIDs(studentIDs, limit, offset)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get achievement references: " + err.Error(), 500, nil))
	}
	
	// Get total count for pagination
	totalCount, err := achievementRefRepo.CountByStudentIDs(studentIDs)
	if err != nil {
		totalCount = 0
	}
	
	// 3. Fetch detail dari MongoDB
	var mongoIDs []string
	for _, ref := range achievementRefs {
		mongoIDs = append(mongoIDs, ref.MongoID)
	}
	
	var achievements []model.Achievement
	if len(mongoIDs) > 0 {
		achievements, err = achievementRepo.FindByIDs(mongoIDs)
		if err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("Failed to get achievement details: " + err.Error(), 500, nil))
		}
	}
	
	// Get student details
	students, err := studentRepo.GetStudentsByAdvisorID(userID)
	if err != nil {
		students = []model.Student{} // Continue even if student details fail
	}
	
	// Calculate pagination
	totalPages := (totalCount + limit - 1) / limit
	
	// 4. Return list dengan pagination
	return c.JSON(utils.SuccessResponse("Advisees achievements retrieved", 200, fiber.Map{
		"achievements": achievements,
		"students": students,
		"pagination": fiber.Map{
			"page": page,
			"limit": limit,
			"total": totalCount,
			"total_pages": totalPages,
		},
		"summary": fiber.Map{
			"total_advisees": len(studentIDs),
			"total_achievements": totalCount,
		},
	}))
}

// ====================================================================
// NOTIFICATION FUNCTIONS (Integrated with Achievement Service)
// ====================================================================

// Notification struct
type Notification struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	RecipientID string            `json:"recipient_id" bson:"recipient_id"`
	SenderID    string            `json:"sender_id" bson:"sender_id"`
	Type        string            `json:"type" bson:"type"`
	Title       string            `json:"title" bson:"title"`
	Message     string            `json:"message" bson:"message"`
	Data        interface{}       `json:"data" bson:"data"`
	IsRead      bool              `json:"is_read" bson:"is_read"`
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	ReadAt      *time.Time        `json:"read_at,omitempty" bson:"read_at,omitempty"`
}

// Get notification collection
func getNotificationCollection() *mongo.Collection {
	return database.MongoClient.Database("achievement_db").Collection("notifications")
}

// Create notification (internal function)
func CreateNotification(recipientID, senderID, notificationType, title, message string, data interface{}) error {
	collection := getNotificationCollection()
	
	notification := &Notification{
		RecipientID: recipientID,
		SenderID:    senderID,
		Type:        notificationType,
		Title:       title,
		Message:     message,
		Data:        data,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, notification)
	return err
}

// Create notification for achievement submission
func CreateAchievementSubmissionNotification(advisorID, studentID, achievementID, achievementTitle string) error {
	title := "New Achievement Submission"
	message := fmt.Sprintf("A student has submitted an achievement '%s' for your verification", achievementTitle)
	
	data := fiber.Map{
		"achievement_id": achievementID,
		"student_id":     studentID,
		"action_type":    "verify_achievement",
	}

	return CreateNotification(advisorID, studentID, "achievement_submitted", title, message, data)
}

// Create notification for achievement verified
func CreateAchievementVerifiedNotification(studentID, advisorID, achievementID, achievementTitle string) error {
	title := "Achievement Verified"
	message := fmt.Sprintf("Your achievement '%s' has been verified by your advisor", achievementTitle)
	
	data := fiber.Map{
		"achievement_id": achievementID,
		"advisor_id":     advisorID,
		"action_type":    "view_achievement",
	}

	return CreateNotification(studentID, advisorID, "achievement_verified", title, message, data)
}

// Create notification for achievement rejected
func CreateAchievementRejectedNotification(studentID, advisorID, achievementID, achievementTitle, reason string) error {
	title := "Achievement Rejected"
	message := fmt.Sprintf("Your achievement '%s' has been rejected. Reason: %s", achievementTitle, reason)
	
	data := fiber.Map{
		"achievement_id": achievementID,
		"advisor_id":     advisorID,
		"reason":         reason,
		"action_type":    "edit_achievement",
	}

	return CreateNotification(studentID, advisorID, "achievement_rejected", title, message, data)
}

// Get notifications for user
func GetUserNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	limit := c.QueryInt("limit", 20)
	
	collection := getNotificationCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"recipient_id": userID}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}}).SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get notifications", 500, nil))
	}
	defer cursor.Close(ctx)

	var notifications []Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to parse notifications", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Notifications retrieved", 200, notifications))
}

// Mark notification as read
func MarkNotificationAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	
	objectID, err := primitive.ObjectIDFromHex(notificationID)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Invalid notification ID", 400, nil))
	}

	collection := getNotificationCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
			"read_at": now,
		},
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to mark notification as read", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Notification marked as read", 200, nil))
}