package service

import (
	"fmt"
	"uas_backend/app/model"
	"uas_backend/app/repository"
	"uas_backend/app/utils"

	"github.com/gofiber/fiber/v2"
)

var notificationRepo *repository.NotificationRepository

func SetNotificationRepo(repo *repository.NotificationRepository) {
	notificationRepo = repo
}

// Create notification (internal function)
func CreateNotification(recipientID, senderID, notificationType, title, message string, data interface{}) error {
	notification := &model.Notification{
		RecipientID: recipientID,
		SenderID:    senderID,
		Type:        notificationType,
		Title:       title,
		Message:     message,
		Data:        data,
	}

	_, err := notificationRepo.Create(notification)
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

// Get notifications for user
func GetUserNotifications(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	
	// Get limit from query parameter (default 20)
	limit := c.QueryInt("limit", 20)
	
	notifications, err := notificationRepo.FindByRecipientID(userID, limit)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to get notifications", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Notifications retrieved", 200, notifications))
}

// Mark notification as read
func MarkNotificationAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	
	err := notificationRepo.MarkAsRead(notificationID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse("Failed to mark notification as read", 500, nil))
	}

	return c.JSON(utils.SuccessResponse("Notification marked as read", 200, nil))
}