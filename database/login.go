package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func loginUser(c *fiber.Ctx) error {
	var email, password string

	if c.Is("json") {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		email = input.Email
		password = input.Password
	} else {
		email = c.FormValue("email")
		password = c.FormValue("password")
	}

	

	// ค้นหาผู้ใช้จากอีเมล
	var user userdatabase
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := userCol.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	//อัปเดต last_active
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, _ = userCol.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"last_active": time.Now()}})

	// ตั้งค่า cookie
	c.Cookie(&fiber.Cookie{
		Name:     "user_email",
		Value:    email,
		Expires: time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})


	user.Password = "" // ไม่คืน password
	return c.JSON(user)
}
