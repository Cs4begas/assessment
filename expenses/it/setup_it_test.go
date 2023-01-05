//go:build integration
// +build integration

package expenses_it_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/cs4begas/assessment/expenses"
	"github.com/cs4begas/assessment/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const serverPort = 2565

func startServer() *http.Server {
	router := gin.Default()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: router,
	}
	go func(c *gin.Engine) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			fmt.Printf("Failed to connect to database: %v", err)
			log.Fatal(err)
		}
		expenses.InitDB(db)
		h := expenses.NewApplication(db)

		c.Use(middleware.AuthMiddleware())
		c.POST("/expenses", h.CreateExpense)
		c.GET("/expenses/:id", h.GetExpenseById)
		c.PUT("/expenses/:id", h.UpdateExpenseById)
		c.GET("/expenses", h.GetAllExpenses)

		fmt.Printf("Server is running at port %d", serverPort)
		c.Run(fmt.Sprintf(":%d", serverPort))
	}(router)

	return srv
}

func connectServer() {
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func shutdownServer(srv *http.Server, t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)
}

func callRequest(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	resp, err := client.Do(req)
	return resp, err
}
