package main

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Admin secret key — ต้องตั้งใน environment variable
var adminSecretKey = os.Getenv("ADMIN_SECRET_KEY")

// PromoteToAdmin — ยกระดับผู้ใช้เป็น admin
// POST /api/admin/promote?secret=YOUR_SECRET_KEY
// body: { "email": "user@example.com" }
func PromoteToAdmin(c *fiber.Ctx) error {
	secret := c.Query("secret")
	if secret == "" || secret != adminSecretKey {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Invalid or missing secret key"})
	}

	type req struct {
		Email string `json:"email"`
	}
	var r req
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if r.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userCol.UpdateOne(ctx, bson.M{"email": r.Email}, bson.M{"$set": bson.M{"role": "admin"}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating user"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"message": "User promoted to admin successfully",
		"email":   r.Email,
	})
}

// CreateAdminUser — สมัครสมาชิก admin (ต้องมี secret key)
// POST /api/admin/register?secret=YOUR_SECRET_KEY
// body: { "email": "admin@example.com", "name": "Admin Name", "password": "secret123" }
func CreateAdminUser(c *fiber.Ctx) error {
	secretKey := c.Query("ADMIN_SECRET_KEY")
	expectedKey := os.Getenv("ADMIN_SECRET_KEY")
	if secretKey != expectedKey || secretKey == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or missing secret key")
	}

	var user userdatabase
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	user.Password = string(hash)
	user.Role = "admin" // ตั้ง role เป็น admin โดยตรง

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := userCol.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error inserting admin user"})
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	user.Password = "" // Do not return password
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin user created successfully",
		"user":    user,
	})
}
