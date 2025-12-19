package service

import (
	"uas_backend/app/model"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

// 5.5 Students & Lecturers

// GetAllStudents godoc
// @Summary Get All Students
// @Description Get list of all students with pagination
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response "Students retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/students [get]
func GetAllStudents(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}

	students, err := studentRepo.FindAllWithPagination(limit, offset)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get students", 500, nil))
	}

	total, err := studentRepo.Count()
	if err != nil {
		total = 0
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(utils.SuccessResponse("Students retrieved", 200, fiber.Map{
		"students": students,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}))
}

// GetStudentByID godoc
// @Summary Get Student by ID
// @Description Get student details by ID
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Success 200 {object} utils.Response "Student retrieved successfully"
// @Failure 404 {object} utils.Response "Student not found"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/students/{id} [get]
func GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")

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

	return c.JSON(utils.SuccessResponse("Student retrieved", 200, student))
}

// GetStudentAchievements godoc
// @Summary Get Student Achievements
// @Description Get all achievements for a specific student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response "Achievements retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/students/{id}/achievements [get]
func GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if achievementRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Achievement repository not initialized", 500, nil))
	}

	// Get student's user_id
	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}

	student, err := studentRepo.FindByID(studentID)
	if err != nil || student == nil {
		return c.Status(404).JSON(utils.ErrorResponse("Student not found", 404, nil))
	}

	// Get achievements by user_id
	filter := map[string]interface{}{
		"student_id": student.UserID,
		"is_deleted": map[string]interface{}{"$ne": true},
	}

	achievements, total, err := achievementRepo.FindAllWithPagination(filter, page, limit, "created_at", "desc")
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get achievements", 500, nil))
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(utils.SuccessResponse("Student achievements retrieved", 200, fiber.Map{
		"student_id":   studentID,
		"achievements": achievements,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}))
}

// UpdateStudentAdvisor godoc
// @Summary Update Student Advisor
// @Description Update advisor for a student
// @Tags Students
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Student ID"
// @Param request body model.UpdateAdvisorRequest true "Advisor data"
// @Success 200 {object} utils.Response "Advisor updated successfully"
// @Failure 400 {object} utils.Response "Bad request"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/students/{id}/advisor [put]
func UpdateStudentAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}

	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}

	err := studentRepo.UpdateAdvisor(id, req.AdvisorID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to update advisor", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Student advisor updated", 200, fiber.Map{
		"student_id": id,
		"advisor_id": req.AdvisorID,
	}))
}

// GetAllLecturers godoc
// @Summary Get All Lecturers
// @Description Get list of all lecturers with pagination
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.Response "Lecturers retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/lecturers [get]
func GetAllLecturers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}

	lecturers, err := studentRepo.FindAllLecturers(limit, offset)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get lecturers", 500, nil))
	}

	total, err := studentRepo.CountLecturers()
	if err != nil {
		total = 0
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(utils.SuccessResponse("Lecturers retrieved", 200, fiber.Map{
		"lecturers": lecturers,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}))
}

// GetLecturerAdvisees godoc
// @Summary Get Lecturer Advisees
// @Description Get all students advised by a specific lecturer
// @Tags Lecturers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Lecturer ID"
// @Success 200 {object} utils.Response "Advisees retrieved successfully"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /v1/lecturers/{id}/advisees [get]
func GetLecturerAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")

	if studentRepo == nil {
		return c.Status(500).JSON(utils.ErrorResponse("Student repository not initialized", 500, nil))
	}

	// Get lecturer's user_id
	lecturer, err := studentRepo.FindLecturerByID(lecturerID)
	if err != nil || lecturer == nil {
		return c.Status(404).JSON(utils.ErrorResponse("Lecturer not found", 404, nil))
	}

	// Get advisees
	advisees, err := studentRepo.GetStudentsByAdvisorID(lecturer.UserID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get advisees", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Lecturer advisees retrieved", 200, fiber.Map{
		"lecturer_id": lecturerID,
		"advisees":    advisees,
		"total":       len(advisees),
	}))
}