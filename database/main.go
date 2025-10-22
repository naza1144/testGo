package main

import (
	"context"
	"log"
	"os"
	"time"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userdatabase struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password,omitempty" json:"password,omitempty"`
}

var userCol *mongo.Collection

func loggingMiddleware(c *fiber.Ctx) error {
  // Start timer
  start := time.Now()

  // Process request
  err := c.Next()

  // Calculate processing time
  duration := time.Since(start)

  // Log the information
  fmt.Printf("Request URL: %s - Method: %s - Duration: %s\n", c.OriginalURL(), c.Method(), duration)

  return err
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
	engine := html.New("D:/code all/testGo/template", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware
	app.Use(loggingMiddleware)

	// ปิดการเก็บแคชสำหรับไฟล์ static ระหว่าง dev
	app.Use("/static", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		return c.Next()
	})
	
    app.Static("/static", "./static")

	// Routes
	app.Post("/api/register", registerUser)
	app.Post("/api/login", loginUser)
	app.Get("/api/users", allUsers)

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
