package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/apache/arrow/go/v12/arrow"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "github.com/dgrijalva/jwt-go"

)

const (
	// sqlQueryData= `SELECT * FROM test_profile WHERE email = ?`
	sqlQueryData = `SELECT name, email, idcard FROM test_profile WHERE email = ?`
	sqlInsertdataprofile= `INSERT INTO test_profile (name, email, idcard) VALUES (?, ?, ?)`
)

//Create Token
type CreateToken struct {
	userID int `json:"userID"`
	Name string `json:"name"` 
	Role string `json:"role"`
	Token string `json:"token"`
}

//AuthLogin

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

var SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))
func (repo *CreateToken) CreateToken(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(), // Token berlaku selama 1 jam
		"email": payload["email"],               // Ambil email dari parameter payload
        "name":  payload["name"],
		"roleID": payload["roleID"], 

	}
	if len(SECRET_KEY) == 0 {
		// log.Fatal("JWT_SECRET is not set")
		log.Fatalf("JWT_SECRET is not set. Current value: %s", os.Getenv("SECRET_KEY"))

	}

	
	// Tambahkan payload ke klaim token
	for key, value := range payload {
		claims[key] = value
	}

	// Buat token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SECRET_KEY)
}

func (repo *CreateToken) Create(c echo.Context) error {
	// Terima dan bind payload dari request body
	payload := make(map[string]interface{})
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"status":  "error",
			"message": "Invalid request body, reqruired name, email, roleID",
		})
	}

	//Create Payload
	name, _ := payload["name"].(string)
	email, _ := payload["email"].(string)
	roleID, _ := payload["roleID"].(string)
	
	// Buat token JWT
	token, err := repo.CreateToken(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "error",
			"message": "Failed to create token",
		})
	}

	// Menambahkan custom header pada response
	c.Response().Header().Set("X-Custom-Header", "CustomHeaderValue")
	c.Response().Header().Set("Authorization", "Bearer "+token)

	// Kembalikan token ke client
	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
		"status": "Success",
		"name": name,
		"email": email,
		"roleID": roleID,
	})
}

//Function Auth Login


func main() {
	e := echo.New()
	godotenv.Load()
	
	
	//JWT_SECRET
	SECRET_KEY = []byte(os.Getenv("JWT_SECRET"))
	if len(SECRET_KEY) == 0 {
		log.Fatal("JWT_SECRET is not set in .env")
	}

	// Middleware logging
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())


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
	repo := &CreateToken{}

	updateIDCard := &UpdateIDCard{
		db: db,
	}
	GetUser := &GetUser{
		db: db,
	}

	//Endpoint List
	e.GET("/get-users", GetUser.Select)
	e.PATCH("/update-idcard", updateIDCard.Update)
	e.POST("/api/create-token", repo.Create)
	e.Logger.Fatal(e.Start(":8080"))
}


// 1. with body , user email dll
// 2. Handler menghandle apa saja yang ada di body
// 3. Create to Claims & output 