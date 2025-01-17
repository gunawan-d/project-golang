package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"  // Import middleware dengan benar

)

type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Timestamp string `json:"timestamp"`

} 

type RefreshServer struct {
	refreshQueue chan int64
	isRefreshing bool
	lastRefreshTime int64
}

func NewRefreshServer() *RefreshServer {
	return &RefreshServer{
		refreshQueue: make(chan int64, 2),
		isRefreshing: false,
		lastRefreshTime: 2,
	}
}

func (s *RefreshServer) StartProcessing() {
	go func () {
		for timestamp :=range s.refreshQueue {
			s.isRefreshing = true
			fmt.Println("Start Refresh proccess at Unix Time:", timestamp)

			time.Sleep(10 * time.Second)

			s.isRefreshing = false
			fmt.Println("Refresh Completed at UNix time:", time.Now().Unix())
		}
	}()
}

func (s *RefreshServer) HandleRefresh(c echo.Context) error {
	currentUnix := time.Now().Unix()

	if currentUnix <= s.lastRefreshTime {
		return c.JSON(http.StatusTooManyRequests, map[string]interface{}{
			"error": "Request to frequent, n please wait",
			"current_time": currentUnix,
			"last_refresh": s.lastRefreshTime,
		})
	}

	select {
	case s.refreshQueue <- currentUnix:
		Response := Response{
			Status: "Created",
			Message: fmt.Sprint("refresh requesd queued at:", currentUnix),
			Timestamp: fmt.Sprint(currentUnix),
		}

		return c.JSON(http.StatusCreated, Response)
	default:
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"error": "refresh queue is full",
			"timestamp": currentUnix,
		})
	}
}

func main() {
	//Create echo isntance
	e :=echo.New()

	// Middeware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	//intialize Server
	server := NewRefreshServer()
	server.StartProcessing()

	//Routers
	e.POST("/refresh", server.HandleRefresh)

	//Healthcheck nanti dulu

	//initial refresh
	fmt.Println("excute intial refresh")
	initialTime := time.Now().Unix()
	server.refreshQueue <- initialTime
	fmt.Println("Initial Refresh Queued at UNIX time:", initialTime)

	//start server
	fmt.Println("Server Running on Port 8080.....")
	e.Logger.Fatal(e.Start(":8080"))


}