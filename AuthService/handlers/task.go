package handlers

import (
	"time"
	"github.com/appu900/authservice/database"
	"github.com/appu900/authservice/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HandleCreateTask(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(primitive.ObjectID)
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	now := time.Now()
	task := models.Task{
		ID:          primitive.NewObjectID(),
		UserID:      userId,
		Title:       input.Title,
		Description: input.Description,
		Status:      "pending",
		DueDate:     now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	taskCollection := database.GetCollection("tasks")
	_, err := taskCollection.InsertOne(c.Context(), task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while creating the task",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task created successfully",
		"task":    task,
	})
}

func HandleGetAlltasksOfUser(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(primitive.ObjectID)

	status := c.Query("status")
	sortBy := c.Query("sort_by", "craeted_at")
	order := c.Query("order", "desc")

	filter := bson.M{"user_id": userId}
	if status != "" {
		filter["status"] = status
	}

	sortDirection := 1
	if order == "desc" {
		sortDirection = -1
	}

	opts := options.Find().SetSort(bson.D{{Key: sortBy, Value: sortDirection}})

	taskCollection := database.GetCollection("tasks")
	cursor, err := taskCollection.Find(c.Context(), filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while fetching tasks",
		})
	}

	defer cursor.Close(c.Context())
	var tasks []models.Task
	if err = cursor.All(c.Context(), &tasks); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while fetching tasks",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tasks": tasks,
		"count": len(tasks),
	})
}

func HandleUpdateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(primitive.ObjectID)
	taskID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Task ID",
		})
	}

	var Input struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&Input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	validStatus := map[string]bool{
		"pending":   true,
		"completed": true,
	}

	if !validStatus[Input.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status",
		})
	}

	taskCollection := database.GetCollection("tasks")
	result, err := taskCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": taskID, "user_id": userID},
		bson.M{"$set": bson.M{
			"status":     "completed",
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while updating the task",
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task not found or Unauthorized",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Task updated successfully",
	})
}




