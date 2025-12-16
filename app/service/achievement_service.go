package service

import (
	"fmt"
	"time"
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

// GET /api/v1/achievements - List (filtered by role)
func GetAllAchievements(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	var filter bson.M
	
	// Role-based filtering
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" { // Mahasiswa role ID
		// Mahasiswa hanya bisa lihat achievement sendiri
		filter = bson.M{"student_id": userID}
	} else {
		// Admin/Dosen bisa lihat semua
		filter = bson.M{}
	}
	
	achievements, err := achievementRepo.FindAll(filter)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get achievements", 500, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Achievements retrieved", 200, achievements))
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

// POST /api/v1/achievements - Create (Mahasiswa) - FR-003 Implementation
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
	
	// TODO: Implement achievement update
	return c.JSON(utils.SuccessResponse("Achievement updated", 200, fiber.Map{"id": id}))
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
	
	// 2. Check ownership - mahasiswa hanya bisa delete achievement sendiri
	if achievement.StudentID != userID {
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

// POST /api/v1/achievements/:id/submit - Submit for verification (FR-004)
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
	
	// 2. Check ownership - mahasiswa hanya bisa submit achievement sendiri
	if achievement.StudentID != userID {
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

// POST /api/v1/achievements/:id/verify - Verify (Dosen Wali)
func VerifyAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.VerifyAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// TODO: Implement achievement verification
	return c.JSON(utils.SuccessResponse("Achievement verified", 200, fiber.Map{"id": id}))
}

// POST /api/v1/achievements/:id/reject - Reject (Dosen Wali)
func RejectAchievement(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.RejectAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// TODO: Implement achievement rejection
	return c.JSON(utils.SuccessResponse("Achievement rejected", 200, fiber.Map{"id": id}))
}

// GET /api/v1/achievements/:id/history - Status history
func GetAchievementHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement get achievement history
	return c.JSON(utils.SuccessResponse("Achievement history retrieved", 200, fiber.Map{"id": id}))
}

// POST /api/v1/achievements/:id/attachments - Upload files
func UploadAchievementAttachment(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement file upload
	return c.JSON(utils.SuccessResponse("Attachment uploaded", 200, fiber.Map{"id": id}))
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