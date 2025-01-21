package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/apache/arrow/go/v12/arrow"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	// sqlQueryData= `SELECT * FROM test_profile WHERE email = ?`
	sqlQueryData = `SELECT name, email, idcard FROM test_profile WHERE email = ?`
	sqlInsertdataprofile= `INSERT INTO test_profile (name, email, idcard) VALUES (?, ?, ?)`
)

type GetUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	IDCard  int    `json:"idcard"`
	db      *sqlx.DB
}

type getMessage struct {
	Status  string `json:"status"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Email   string `json:"email"`
	IDCard  int    `json:"idcard"`
	db      *sqlx.DB
}

type UpdateIDCard struct {
	IDCard string `json:"idcard"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	db    *sqlx.DB
}

type responseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//Update IDCard
func (repo *UpdateIDCard) Update(c echo.Context) error {
	req := new(UpdateIDCard)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, responseMessage{
			Status:  "error",
			Message: "Invalid request body",
		})
	}

	if len(req.IDCard) != 16 {
		return c.JSON(http.StatusBadRequest, responseMessage{
			Status:  "error",
			Message: "Invalid ID KTP, must be 16 characters",
		})
	}

	// query := sqlInsertdataprofile
	res, err := repo.db.Exec(sqlInsertdataprofile, req.Name, req.Email, req.IDCard)
	if err != nil {
		log.Printf("Failed to insert data into database: %v", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{
			Status:  "error",
			Message: "Failed to insert data",
		})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, responseMessage{
			Status:  "error",
			Message: "No user found with the specified ID",
		})
	}

	return c.JSON(http.StatusOK, responseMessage{
		Status:  "success",
		Message: "ID card updated successfully",
	})
}


//Get User with param email
func (repo *GetUser) Select(c echo.Context) error {
	response := new(GetUser)

	email := c.FormValue("email")
	err := repo.db.QueryRow(sqlQueryData, email).Scan(&response.Name, &response.Email, &response.IDCard)

	log.Printf("Failed to insert data into database: %v", err)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, getMessage{
			Status:  "error",
			Message: "Failed to Get data",
		})
	}

	return c.JSON(http.StatusOK, response)

}

func main() {
	e := echo.New()
	godotenv.Load()

	// Middleware logging
	e.Use(middleware.Logger())

	fmt.Println("execute start")
	fmt.Println("Server running on port 8080....")
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
	log.Println("Database Connected")


	//Handler List
	updateIDCard := &UpdateIDCard{
		db: db,
	}
	GetUser := &GetUser{
		db: db,
	}

	//Endpoint List
	e.GET("/get-users", GetUser.Select)
	e.PATCH("/update-idcard", updateIDCard.Update)
	e.Logger.Fatal(e.Start(":8080"))
}
