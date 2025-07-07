package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	HOST     = "dpg-d1in3u0dl3ps73d35oqg-a.oregon-postgres.render.com"
	PORT     = "5432"
	USERNAME = "api_database_jqur_user"
	PASSWORD = "jx7YxQwWkqOdLAavkYl53O3AFCSGRQzu"
	DBNAME   = "api_database_jqur"
)

func GetPsqlInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		HOST, PORT, USERNAME, PASSWORD, DBNAME)
}

func CreateDbObject() error {
	var err error
	DB, err = sql.Open("postgres", GetPsqlInfo())
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to DB: %w", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	DB.SetConnMaxLifetime(1 * time.Hour)

	return nil
}

type FetchAllUsersOutput struct {
	UserID int64
	Name   string
	Email  string
}

func FetchAllUsers() ([]FetchAllUsersOutput, error) {
	query := "SELECT userid, name, email FROM users"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []FetchAllUsersOutput
	for rows.Next() {
		var u FetchAllUsersOutput
		err := rows.Scan(&u.UserID, &u.Name, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func main() {
	if err := CreateDbObject(); err != nil {
		log.Fatal("❌ DB Connection Failed:", err)
	}

	app := fiber.New()

	// Root route to confirm deployment
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Go app is successfully deployed on Render!")
	})

	// Endpoint to fetch users
	app.Get("/users", func(c *fiber.Ctx) error {
		users, err := FetchAllUsers()
		if err != nil {
			return c.Status(500).SendString("❌ Failed to fetch users: " + err.Error())
		}
		return c.JSON(users)
	})

	// ✅ Register SignIn route
	app.Post("/signin", SignInWeb)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInWeb(c *fiber.Ctx) error {
	signInRequestObject := SignInRequest{}
	if err := c.BodyParser(&signInRequestObject); err != nil {
		c.Status(500)
		return nil
	}
	fmt.Println(signInRequestObject)
	c.JSON("Its WORKING")
	return nil
}
