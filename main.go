package main

import (
	"fmt"
	"log"
	"os"
	

	"github.com/joho/godotenv"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"project-golang/repository"
	"project-golang/rest"
	"project-golang/services"
	_ "github.com/go-sql-driver/mysql"
	
)

func main() {
	godotenv.Load()

	//JWT_SECRET
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DBUser"),
		os.Getenv("DBPass"),
		os.Getenv("DBHost"),
		os.Getenv("DBPort"),
		os.Getenv("DBName"),
	)
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize Repositories
	userRepo := &repository.UserRepository{DB: db}

	// Initialize Services
	authService := &services.AuthService{}

	// Initialize Handlers
	userHandler := &rest.UserHandler{Repo: userRepo}
	authHandler := &rest.AuthHandler{AuthService: authService}

	// Routes
	e.GET("/get-users", userHandler.GetUser)
	e.PATCH("/update-idcard", userHandler.UpdateIDCard)
	e.POST("/api/create-token", authHandler.CreateToken)

	// Start Server
	log.Println("Server running on port 8080...")
	e.Logger.Fatal(e.Start(":8080"))
}
