package service

import (
	"context"
	"fmt"
	"time"
	"uas_backend/app/model"
	"uas_backend/app/utils"
	"uas_backend/database"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetStatistics godoc
// @Summary Get Achievement Statistics
// @Description Get comprehensive achievement statistics with role-based filtering
// @Tags Reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} utils.Response "Statistics retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/reports/statistics [get]
func GetStatistics(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	

	
	// Get query parameters for filtering
	startDate := c.Query("start_date") // YYYY-MM-DD
	endDate := c.Query("end_date")     // YYYY-MM-DD
	
	// Set default period (last 12 months)
	periodEnd := time.Now()
	periodStart := periodEnd.AddDate(-1, 0, 0)
	
	// Parse custom date range if provided
	if startDate != "" {
		if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
			periodStart = parsed
		}
	}
	if endDate != "" {
		if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
			periodEnd = parsed
		}
	}
	
	// Build filter based on role
	var filter bson.M
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" { // Mahasiswa
		// Mahasiswa hanya lihat statistik sendiri
		filter = bson.M{
			"student_id": userID,
			"is_deleted": bson.M{"$ne": true},
			"created_at": bson.M{
				"$gte": periodStart,
				"$lte": periodEnd,
			},
		}
	} else if role == "9f6c3a32-ba48-4f0a-a69e-d89be58a2d8e" { // Dosen
		// Dosen lihat statistik mahasiswa bimbingannya
		if studentRepo == nil {
			return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
		}
		
		studentIDs, err := studentRepo.GetStudentIDsByAdvisorID(userID)
		if err != nil {
			return c.Status(500).JSON(utils.ErrorResponse("Failed to get advisees", 500, nil))
		}
		
		filter = bson.M{
			"student_id": bson.M{"$in": studentIDs},
			"is_deleted": bson.M{"$ne": true},
			"created_at": bson.M{
				"$gte": periodStart,
				"$lte": periodEnd,
			},
		}
	} else if role == "fd796792-3c30-4e34-b2c8-fa2f93d201e7" { // Admin
		// Admin lihat semua statistik
		filter = bson.M{
			"is_deleted": bson.M{"$ne": true},
			"created_at": bson.M{
				"$gte": periodStart,
				"$lte": periodEnd,
			},
		}
	} else {
		// Default: hanya data sendiri (fallback untuk role tidak dikenal)
		filter = bson.M{
			"student_id": userID,
			"is_deleted": bson.M{"$ne": true},
			"created_at": bson.M{
				"$gte": periodStart,
				"$lte": periodEnd,
			},
		}
	}
	
	// Generate statistics
	stats, err := generateAchievementStatistics(filter, periodStart, periodEnd)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to generate statistics", 500, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Achievement statistics retrieved", 200, stats))
}

// Generate comprehensive achievement statistics
func generateAchievementStatistics(filter bson.M, periodStart, periodEnd time.Time) (*model.AchievementStatistics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	collection := database.MongoClient.Database("achievement_db").Collection("achievements")
	
	stats := &model.AchievementStatistics{
		ByCategory:         make(map[string]int),
		ByPeriod:          make(map[string]int),
		TopStudents:       []model.TopStudent{},
		ByCompetitionLevel: make(map[string]int),
		PeriodStart:       periodStart,
		PeriodEnd:         periodEnd,
	}
	
	// 1. Basic counts by status
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": "$status",
			"count": bson.M{"$sum": 1},
			"total_points": bson.M{"$sum": "$points"},
		}},
	}
	
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	// Process status counts
	for cursor.Next(ctx) {
		var result struct {
			ID          string `bson:"_id"`
			Count       int    `bson:"count"`
			TotalPoints int    `bson:"total_points"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		
		stats.TotalAchievements += result.Count
		stats.TotalPoints += result.TotalPoints
		
		switch result.ID {
		case "verified":
			stats.VerifiedCount = result.Count
		case "submitted":
			stats.PendingCount = result.Count
		case "rejected":
			stats.RejectedCount = result.Count
		case "draft":
			stats.DraftCount = result.Count
		}
	}
	
	// 2. Count by category
	pipeline = []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": "$category",
			"count": bson.M{"$sum": 1},
		}},
	}
	
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err == nil {
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int    `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				stats.ByCategory[result.ID] = result.Count
			}
		}
		cursor.Close(ctx)
	}
	
	// 3. Count by period (monthly)
	pipeline = []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$created_at"},
				"month": bson.M{"$month": "$created_at"},
			},
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"_id.year": 1, "_id.month": 1}},
	}
	
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err == nil {
		for cursor.Next(ctx) {
			var result struct {
				ID struct {
					Year  int `bson:"year"`
					Month int `bson:"month"`
				} `bson:"_id"`
				Count int `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				period := fmt.Sprintf("%d-%02d", result.ID.Year, result.ID.Month)
				stats.ByPeriod[period] = result.Count
			}
		}
		cursor.Close(ctx)
	}
	
	// 4. Competition level distribution (based on points)
	stats.ByCompetitionLevel["Lokal (1-25 poin)"] = 0
	stats.ByCompetitionLevel["Regional (26-50 poin)"] = 0
	stats.ByCompetitionLevel["Nasional (51-75 poin)"] = 0
	stats.ByCompetitionLevel["Internasional (76+ poin)"] = 0
	
	pipeline = []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": bson.M{
				"$switch": bson.M{
					"branches": []bson.M{
						{"case": bson.M{"$lte": []interface{}{"$points", 25}}, "then": "Lokal (1-25 poin)"},
						{"case": bson.M{"$lte": []interface{}{"$points", 50}}, "then": "Regional (26-50 poin)"},
						{"case": bson.M{"$lte": []interface{}{"$points", 75}}, "then": "Nasional (51-75 poin)"},
					},
					"default": "Internasional (76+ poin)",
				},
			},
			"count": bson.M{"$sum": 1},
		}},
	}
	
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err == nil {
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int    `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				stats.ByCompetitionLevel[result.ID] = result.Count
			}
		}
		cursor.Close(ctx)
	}
	
	// Calculate average points
	if stats.TotalAchievements > 0 {
		stats.AveragePoints = float64(stats.TotalPoints) / float64(stats.TotalAchievements)
	}
	
	return stats, nil
}

// GET /api/v1/reports/student/:id
func GetStudentReport(c *fiber.Ctx) error {
	id := c.Params("id")
	
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	
	// Check access permissions
	if role == "f464ceb1-5481-49cf-99f0-d8f2d66f4506" { // Mahasiswa
		// Mahasiswa hanya bisa lihat report sendiri
		if studentRepo == nil {
			return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
		}
		
		student, err := studentRepo.FindByUserID(userID)
		if err != nil || student == nil {
			return c.Status(404).JSON(utils.ErrorResponse("Student not found", 404, nil))
		}
		
		if student.ID != id {
			return c.Status(403).JSON(utils.ErrorResponse("Access denied - can only view own report", 403, nil))
		}
	}
	
	// Get student details
	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}
	
	student, err := studentRepo.FindByID(id)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get student", 500, nil))
	}
	if student == nil {
		return c.Status(404).JSON(utils.ErrorResponse("Student not found", 404, nil))
	}
	
	// Generate student statistics
	filter := bson.M{
		"student_id": student.UserID,
		"is_deleted": bson.M{"$ne": true},
	}
	
	stats, err := generateStudentStatistics(filter, student)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to generate student report", 500, nil))
	}
	
	return c.JSON(utils.SuccessResponse("Student report retrieved", 200, stats))
}
// Generate student-specific statistics
func generateStudentStatistics(filter bson.M, student *model.Student) (*model.StudentReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	collection := database.MongoClient.Database("achievement_db").Collection("achievements")
	
	report := &model.StudentReport{
		StudentID:   student.ID,
		StudentName: student.Name,
		NIM:         student.NIM,
		Program:     student.Program,
		Semester:    student.Semester,
	}
	
	// Basic counts by status
	pipeline := []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": "$status",
			"count": bson.M{"$sum": 1},
			"total_points": bson.M{"$sum": "$points"},
		}},
	}
	
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	// Process status counts
	for cursor.Next(ctx) {
		var result struct {
			ID          string `bson:"_id"`
			Count       int    `bson:"count"`
			TotalPoints int    `bson:"total_points"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		
		report.TotalAchievements += result.Count
		report.TotalPoints += result.TotalPoints
		
		switch result.ID {
		case "verified":
			report.VerifiedCount = result.Count
		case "submitted":
			report.PendingCount = result.Count
		case "rejected":
			report.RejectedCount = result.Count
		case "draft":
			report.DraftCount = result.Count
		}
	}
	
	// Calculate average points
	if report.TotalAchievements > 0 {
		report.AveragePoints = float64(report.TotalPoints) / float64(report.TotalAchievements)
	}
	
	// Get achievements by category
	report.ByCategory = make(map[string]int)
	pipeline = []bson.M{
		{"$match": filter},
		{"$group": bson.M{
			"_id": "$category",
			"count": bson.M{"$sum": 1},
		}},
	}
	
	cursor, err = collection.Aggregate(ctx, pipeline)
	if err == nil {
		for cursor.Next(ctx) {
			var result struct {
				ID    string `bson:"_id"`
				Count int    `bson:"count"`
			}
			if err := cursor.Decode(&result); err == nil {
				report.ByCategory[result.ID] = result.Count
			}
		}
		cursor.Close(ctx)
	}
	
	return report, nil
}