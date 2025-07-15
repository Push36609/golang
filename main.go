// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// const (
// 	HOST     = "dpg-d1in3u0dl3ps73d35oqg-a.oregon-postgres.render.com"
// 	PORT     = "5432"
// 	USERNAME = "api_database_jqur_user"
// 	PASSWORD = "jx7YxQwWkqOdLAavkYl53O3AFCSGRQzu"
// 	DBNAME   = "api_database_jqur"
// )

// func GetPsqlInfo() string {
// 	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
// 		HOST, PORT, USERNAME, PASSWORD, DBNAME)
// }

// func CreateDbObject() error {
// 	var err error
// 	DB, err = sql.Open("postgres", GetPsqlInfo())
// 	if err != nil {
// 		return fmt.Errorf("error opening DB: %w", err)
// 	}

// 	err = DB.Ping()
// 	if err != nil {
// 		return fmt.Errorf("error connecting to DB: %w", err)
// 	}

// 	fmt.Println("✅ Connected to PostgreSQL")

// 	DB.SetMaxOpenConns(25)
// 	DB.SetMaxIdleConns(25)
// 	DB.SetConnMaxIdleTime(10 * time.Minute)
// 	DB.SetConnMaxLifetime(1 * time.Hour)

// 	return nil
// }

// type FetchAllUsersOutput struct {
// 	UserID int64
// 	Name   string
// 	Email  string
// }

// func FetchAllUsers() ([]FetchAllUsersOutput, error) {
// 	query := "SELECT userid, name, email FROM users"

// 	rows, err := DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []FetchAllUsersOutput
// 	for rows.Next() {
// 		var u FetchAllUsersOutput
// 		err := rows.Scan(&u.UserID, &u.Name, &u.Email)
// 		if err != nil {
// 			return nil, err
// 		}
// 		users = append(users, u)
// 	}

// 	return users, nil
// }

// // ✅ FIXED: Query from "users" instead of "user"
// func FetchUserIDFromEmailID(email string) (int, error) {
// 	query := `SELECT userid FROM users WHERE email = $1`

// 	var userID int
// 	err := DB.QueryRow(query, email).Scan(&userID)
// 	if err == sql.ErrNoRows {
// 		return 0, fmt.Errorf("no user found with email: %s", email)
// 	}
// 	if err != nil {
// 		return 0, err
// 	}

// 	return userID, nil
// }
// func FetchPasswordByUserID(userID int) (string, error) {
// 	query := `SELECT password FROM users WHERE userid = $1`

// 	var password string
// 	err := DB.QueryRow(query, userID).Scan(&password)
// 	if err == sql.ErrNoRows {
// 		return "", fmt.Errorf("no user found with userID: %d", userID)
// 	}
// 	if err != nil {
// 		return "", err
// 	}

// 	return password, nil
// }
// func main() {
// 	if err := CreateDbObject(); err != nil {
// 		log.Fatal("❌ DB Connection Failed:", err)
// 	}

// 	app := fiber.New()

// 	// Root route to confirm deployment
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendString("✅ Go app is successfully deployed on Render!")
// 	})

// 	// Endpoint to fetch users
// 	app.Get("/users", func(c *fiber.Ctx) error {
// 		users, err := FetchAllUsers()
// 		if err != nil {
// 			return c.Status(500).SendString("❌ Failed to fetch users: " + err.Error())
// 		}
// 		return c.JSON(users)
// 	})

// 	// ✅ Register SignIn route
// 	app.Post("/signin", SignInWeb)

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "3000"
// 	}

// 	log.Printf("🚀 Server running on port %s", port)
// 	log.Fatal(app.Listen("0.0.0.0:" + port))
// }

// type SignInRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// func SignInWeb(c *fiber.Ctx) error {
// 	signInRequestObject := SignInRequest{}
// 	if err := c.BodyParser(&signInRequestObject); err != nil {
// 		c.Status(500)
// 		return nil
// 	}
// 	fmt.Println(signInRequestObject)
// 	c.JSON("Its WORKING")
// 	return nil
// }

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

func main() {
	if err := connectDB(); err != nil {
		log.Fatal("❌ DB Connection Failed:", err)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Go app is successfully deployed on Render!")
	})

	app.Post("/users", getAllUsers)
	app.Post("/signin", signInHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("🚀 Server running on port %s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func connectDB() error {
	var err error
	psql := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", HOST, PORT, USERNAME, PASSWORD, DBNAME)
	DB, err = sql.Open("postgres", psql)
	if err != nil {
		return err
	}
	if err = DB.Ping(); err != nil {
		return err
	}
	DB.SetConnMaxIdleTime(10 * time.Minute)
	DB.SetConnMaxLifetime(1 * time.Hour)
	fmt.Println("✅ Connected to PostgreSQL")
	return nil
}

// ---------------------- Handlers ----------------------

func signInHandler(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserId   int64  `json:"userId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return badRequest(c, "Invalid request format")
	}

	// Fetch userId and password from DB using email
	var userID int
	var storedPassword string

	err := DB.QueryRow(`
	SELECT u.userid, p.password , u.email
	FROM users u
	JOIN passwords p ON u.userid = p.userid
	WHERE u.email = $1 AND u.userid = $2
`, req.Email, req.UserId).Scan(&userID, &storedPassword)

	if err != nil {
		return notFound(c, "User not found with this email")
	}

	// Match password
	if req.Password != storedPassword {
		return badRequest(c, "Incorrect password")
	}

	// Success response
	return c.JSON(fiber.Map{
		"message": "welcome mr. Pushpendra🎉 login successfully 👍 ",
		// "db name": DBNAME,
		"userId":   userID,
		"email":    req.Email,
		"password": PASSWORD,
	})
}

func getAllUsers(c *fiber.Ctx) error {
	rows, err := DB.Query("SELECT userid, name, email FROM users")
	if err != nil {
		return serverError(c, err.Error())
	}
	defer rows.Close()

	var users []struct {
		UserID int64  `json:"userId"`
		Name   string `json:"name"`
		Email  string `json:"email"`
	}

	for rows.Next() {
		var u struct {
			UserID int64  `json:"userId"`
			Name   string `json:"name"`
			Email  string `json:"email"`
		}
		if err := rows.Scan(&u.UserID, &u.Name, &u.Email); err != nil {
			return serverError(c, err.Error())
		}
		users = append(users, u)
	}

	return c.JSON(users)
}

// ---------------------- DB Helpers ----------------------

func getUserIDByEmail(email string) (int, error) {
	var userID int
	err := DB.QueryRow("SELECT userid FROM users WHERE email = $1", email).Scan(&userID)
	return userID, err
}

func getEmailByUserID(userID int) (string, error) {
	var email string
	err := DB.QueryRow("SELECT email FROM users WHERE userid = $1", userID).Scan(&email)
	return email, err
}

func getPasswordByUserID(userID int) (string, error) {
	var password string
	err := DB.QueryRow("SELECT password FROM passwords WHERE userid = $1", userID).Scan(&password)
	return password, err
}

// ---------------------- Reusable Responses ----------------------

func badRequest(c *fiber.Ctx, msg string) error {
	return c.Status(400).JSON(fiber.Map{"error": msg})
}

func notFound(c *fiber.Ctx, msg string) error {
	return c.Status(404).JSON(fiber.Map{"error": msg})
}

func serverError(c *fiber.Ctx, msg string) error {
	return c.Status(500).JSON(fiber.Map{"error": msg})
}
