package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func allUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// projection: ดึงเฉพาะ name และ email
	projection := bson.M{"name": 1, "email": 1} // _id:0 ไม่เอา _id
	opts := options.Find().SetProjection(projection)

	cursor, err := userCol.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot fetch users"})
	}
	defer cursor.Close(ctx)

	// ใช้ bson.M เพราะเราไม่ต้อง mapping password
	var users []bson.M
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse users"})
	}
	return c.JSON(users)
}