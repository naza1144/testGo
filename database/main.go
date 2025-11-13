package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userdatabase struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password,omitempty" json:"password,omitempty"`
	Role       string             `bson:"role,omitempty" json:"role,omitempty"`
	LastActive time.Time          `bson:"last_active,omitempty" json:"last_active,omitempty"`
}

var userCol *mongo.Collection

func loggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	latency := time.Since(start)
	log.Printf("%s %s %d %s\n", c.Method(), c.Path(), c.Response().StatusCode(), latency)
	return err
}

// middleware ที่อัปเดต last_active หากเจอ cookie user_email
func updateLastActive(c *fiber.Ctx) error {
	email := c.Cookies("user_email")
	if email != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, _ = userCol.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"last_active": time.Now()}})
	}
	return c.Next()
}

func main() {
	// โหลด env
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DB_NAME")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is required")
	}
	if dbName == "" {
		dbName = "testdb"
	}

	// Connect MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// ✅ ตรวจสอบการเชื่อมต่อ
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	} else {
		fmt.Println("✅ Connected to MongoDB successfully!")
	}

	userCol = client.Database(dbName).Collection("userdatabase")

	// Fiber app
	// Use a relative path so the app can find templates regardless of absolute drive/path
	engine := html.New("./template", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(loggingMiddleware)

	// อัปเดต last_active ถ้ามี cookie ของ user
	app.Use(updateLastActive)

	// ปิดการเก็บแคชสำหรับไฟล์ static ระหว่าง dev
	app.Use("/static", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		return c.Next()
	})

	app.Static("/static", "./static")

	// Routes
	app.Get("/admin", AdminOnly, func(c *fiber.Ctx) error {
		return c.Render("admin", fiber.Map{})
	})
	app.Post("/api/register", registerUser)
	app.Post("/api/login", loginUser)
	app.Get("/api/status", userStatus)

	// Admin routes — ต้องมี secret key ใน query string
	app.Post("/api/admin/register", CreateAdminUser) // สมัครสมาชิก admin
	app.Post("/api/admin/promote", PromoteToAdmin)   // ยกระดับ user เป็น admin

	app.Get("/:page", func(c *fiber.Ctx) error {
		page := c.Params("page")
		return c.Render(page, fiber.Map{})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	log.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))

}

// nodemon --exec go run . --signal SIGTERM
