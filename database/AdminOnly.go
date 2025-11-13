package main

import (

	"github.com/gofiber/fiber/v2"
)


func AdminOnly(c *fiber.Ctx) error {
    role := c.Cookies("user_role")
    if role != "admin" {
        // ถ้าเป็น API ให้ตอบเป็น 403 หรือ redirect ไปหน้า login
        return c.Status(fiber.StatusForbidden).SendString("Access denied")
    }
    return c.Next()
}