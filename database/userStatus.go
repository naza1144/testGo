package main

import (
    "context"
    "time"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
)

func userStatus(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		email = c.Cookies("user_email")
	}
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var u userdatabase
	err := userCol.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// กำหนด threshold ว่าเรียกว่า online ถ้าภายใน 2 นาที
	threshold := 2 * time.Minute
	online := time.Since(u.LastActive) <= threshold

	return c.JSON(fiber.Map{
		"email":	u.Email,
		"lastActive": u.LastActive,
		"online":	online,
	})
}