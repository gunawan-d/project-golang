package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// Response adalah struktur untuk repsonse JSON
type Response struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Timestamp string `json:"timestamp"`
}

//REfreshSErver mengelola state dan operasi Refresh
type RefreshServer struct {
	refreshQueue chan int64
	lastRefreshTime int64
	isRefreshing bool
}

func NewRefreshServer() *RefreshServer {
    return &RefreshServer{
        refreshQueue:    make(chan int64, 100),
        isRefreshing:   false,
        lastRefreshTime: 0,
    }
}

// NewRefreshServer membuat instance baru refreshServer
func NewRefreshHandler() *RefreshServer {
	return &RefreshServer{
		refreshQueue: make(chan int64, 100),
		isRefreshing: false,
		lastRefreshTime: 0,
	}
}

func (s *RefreshServer) StartProcessing() {
	go func ()  {
		for timestamp := range s.refreshQueue {
			s.isRefreshing = true
			fmt.Println("Starting Refresh Process at UNIX time:", timestamp)

			time.Sleep(10 * time.Second)

			s.isRefreshing = false
			fmt.Println("Refresh completed at UNIX time:", time.Now().Unix())
		}
	}()
}

func (s *RefreshServer) HandleRefresh(c *gin.Context) {
	currentUnix := time.Now().Unix()

	//Cek Request Dalam detik yang sama
	if currentUnix <= s.lastRefreshTime {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Request to frequest, Please wait",
			"current_time": currentUnix,
			"last_refresh": s.lastRefreshTime,
		})
		return
	}

	//Update waktu refresh terakhir
	s.lastRefreshTime = currentUnix

	//Kirim ke queue dan return response segera
	select {
	case s.refreshQueue <- currentUnix:
		response := Response{
			Status: "CREATED", //Ubah satatus Menjadi CREATED
			Message: fmt.Sprint("Refresh Requeste Queued at:", currentUnix),
			Timestamp: fmt.Sprint(currentUnix),
		}

		c.JSON(http.StatusCreated, response) //status 201 Created
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":"Refresh Queue is Full",
			"timestamp": currentUnix,
		})
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	server := NewRefreshServer()
	server.StartProcessing()

	r := gin.Default()

	//POST /refresh di port 8080
	r.POST("/refresh", server.HandleRefresh) //Endpoint POST /refresh

	//.... Kode Lainya ....

	fmt.Println("Server Running on Port 8080....")
	r.Run(":8080")
}
	
