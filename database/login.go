package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func loginUser(c *fiber.Ctx) error {
	type creds struct {
		Email	string `json:"email"`
		Password string `json:"password"`
	}
	var b creds
	if err := c.BodyParser(&b); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()


	// ค้นหาผู้ใช้จากอีเมล
	var user userdatabase
	err := userCol.FindOne(ctx, bson.M{"email": b.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(b.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	//อัปเดต last_active
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = userCol.UpdateOne(ctx, bson.M{"email": b.Email}, bson.M{"$set": bson.M{"last_active": time.Now()}})

	// ตั้งค่า cookie
	c.Cookie(&fiber.Cookie{
		Name:     "user_email",
		Value:    user.Email,
		HTTPOnly: true,
		Expires: time.Now().Add(24 * time.Hour),
	})
	c.Cookie(&fiber.Cookie{
		Name: "user_role",
		Value: user.Role,
		HTTPOnly: true,
		Expires: time.Now().Add(24 * time.Hour),
	})


	user.Password = "" // ไม่คืน password
	return c.JSON(fiber.Map{
		"ok":   true,
		"role": user.Role,
	})
}
