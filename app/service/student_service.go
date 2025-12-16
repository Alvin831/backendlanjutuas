package service

import (
	"uas_backend/app/model"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

// 5.5 Students & Lecturers

// GET /api/v1/students
func GetAllStudents(c *fiber.Ctx) error {
	// TODO: Implement get all students
	return c.JSON(utils.SuccessResponse("Students retrieved", 200, []model.Student{}))
}

// GET /api/v1/students/:id
func GetStudentByID(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement get student by ID
	return c.JSON(utils.SuccessResponse("Student retrieved", 200, fiber.Map{"id": id}))
}

// GET /api/v1/students/:id/achievements
func GetStudentAchievements(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement get student achievements
	return c.JSON(utils.SuccessResponse("Student achievements retrieved", 200, fiber.Map{"student_id": id}))
}

// PUT /api/v1/students/:id/advisor
func UpdateStudentAdvisor(c *fiber.Ctx) error {
	id := c.Params("id")
	var req model.UpdateAdvisorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse("Bad request", 400, nil))
	}
	
	// TODO: Implement update student advisor
	return c.JSON(utils.SuccessResponse("Student advisor updated", 200, fiber.Map{"student_id": id, "advisor_id": req.AdvisorID}))
}

// GET /api/v1/lecturers
func GetAllLecturers(c *fiber.Ctx) error {
	// TODO: Implement get all lecturers
	return c.JSON(utils.SuccessResponse("Lecturers retrieved", 200, []model.Lecturer{}))
}

// GET /api/v1/lecturers/:id/advisees
func GetLecturerAdvisees(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement get lecturer advisees
	return c.JSON(utils.SuccessResponse("Lecturer advisees retrieved", 200, fiber.Map{"lecturer_id": id}))
}