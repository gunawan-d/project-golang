package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shirou/gopsutil/v3/disk"
)

// const Database

const (
	sqlSelectCurrentDate = `SELECT CURRENT_DATE;`
)

type MYSQLConfig struct {
	DBHost string
	DBName string
	DBUser string
	DBPass string
	DBPort string
	
}

type MYSQLRepository struct {
	Config *MYSQLConfig
	db *sqlx.DB
}

// var (
// 	env conf.Environment
// )


//struct untuk rerepoonse API
type Response struct {
	Status string `json:"status"`
	Message string `json:"message""`
	Timestamp string `json:"timestamp""`
}

//struct untuk Start
type StartProcessing struct {
	HitCount int64
	StopSignal chan bool
	Running bool

}

type diskHealthCheck struct{}


//Inialisasi Proses
func NewProcessing() *StartProcessing {
	return &StartProcessing{
		HitCount: 0,
		StopSignal: make(chan bool, ),
		Running: false,
	}

}

//function
func (repo *StartProcessing) Start(c echo.Context) error{
	repo.Running = false
	// go func () {
	// 	for i := 1; i <= 5; i++ {
	// 		repo.HitCount++
	// 		log.Println("Hit ke :", repo.HitCount)
	// 	}
	// 	repo.StopSignal <- true
	// 	log.Println("Log Proses")
	// }()

	return c.JSON(http.StatusCreated, map[string]string{
		"status":    "success",
		"message":   "Proses start berjalan!",
		"timestamp": time.Now().Format(time.RFC822),
	})
	
}

//func stop
func (repo *StartProcessing) Stop(c echo.Context) error{
	repo.Running = false

	return c.JSON(http.StatusCreated, map[string]string{
		"status": "Success",
		"message": "Proses dihentikan",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

//func Reload
func (repo *StartProcessing) Reload(c echo.Context) error{
	repo.HitCount = 0
	return c.JSON(http.StatusCreated, map[string]string{
		"status": "Success",
		"message": "Nilai Hit telah di reset",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

//func Healtcheck service
// func (repo *StartProcessing) Healthcheck(c echo.Context) error{
// 	repo.Running = false

// 	return c.JSON(http.StatusOK, map[string]string{
// 		"status": "Healtcheck OK",
// 		"timestamp": time.Now().Format(time.RFC3339),
// 	})
// }

//func Healtcheck with database
func (repo *MYSQLRepository) Healthcheck(c echo.Context) error{
	// Membentuk DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		// repo.Config.DBUser,
		// repo.Config.DBPass,
		// repo.Config.DBHost,
		// repo.Config.DBPort,
		// repo.Config.DBName,
		os.Getenv("DBUser"),
		os.Getenv("DBPass"),
		os.Getenv("DBHost"),
		os.Getenv("DBPort"),
		os.Getenv("DBName"),
		
	)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Database connection error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "error",
			"message": "Failed to Connection DB",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}
	defer db.Close()

	var currentDate string
	err = db.QueryRow(sqlSelectCurrentDate).Scan(&currentDate)
	if err != nil {
		return fmt.Errorf("database health check failed: %v", err)

	}

	fmt.Println("Database is healthy, current date:", currentDate)
	return c.JSON(http.StatusOK, map[string]string{
		"status":    "success",
		"message":   "Database is healthy",
		"date":      currentDate,
		"timestamp": time.Now().Format(time.RFC3339),
	})
	

}


func (repo *diskHealthCheck) health(c echo.Context) error {
	usage, err := disk.Usage("/")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error checking disk usage",
		})
	}

	if usage.UsedPercent > 80 {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "critical",
			"message": fmt.Sprintf("Disk usage critical: %.2f%% used", usage.UsedPercent),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": fmt.Sprintf("Disk usage healthy: %.2f%% used", usage.UsedPercent),
	})
}


func main() {
	e := echo.New()
	// Inisialisasi Repository
	repo := &MYSQLRepository{}
	godotenv.Load()
	diskRepo := &diskHealthCheck{}


	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file:", err)
	// }
	// fmt.Println("DB_HOST:", os.Getenv("DBHost"))
	// fmt.Println("DB_USER:", os.Getenv("DBUser"))
	// fmt.Println("DB_PASS:", os.Getenv("DBPass"))
	// fmt.Println("DB_PORT:", os.Getenv("DBPort"))
	// fmt.Println("DB_NAME:", os.Getenv("DBName"))
	
	// envconfig.MustProcess("", &env)
	
	// Middeware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	server := NewProcessing()

	e.POST("/start", server.Start)
	e.POST("/stop", server.Stop)
	e.POST("/reload", server.Reload)
	// e.GET("/healthcheck", server.Healthcheck)
	e.GET("/healthcheck", repo.Healthcheck)
	e.GET("/disk-health", diskRepo.health)



	// Jalankan HealthCheck
	// err := repo.Healthcheck()
	// if err != nil {
	// 	fmt.Println("HealthCheck Error:", err)
	// } else {
	// 	fmt.Println("Database connection is successful!")
	// }

	//initial start
	fmt.Println("execute start")

	fmt.Println("Server running on port 8080....")
	e.Logger.Fatal(e.Start(":8080"))
}