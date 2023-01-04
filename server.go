package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cs4begas/assessment/expenses"
	"github.com/cs4begas/assessment/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Please use server.go for main file")
	db := expenses.InitDB()

	h := expenses.NewApplication(db)
	router := SetupRouter(h)
	os_port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    ":" + os_port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("start at port:", os_port)

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown :", err)
	}
}

func SetupRouter(h *expenses.Handler) *gin.Engine {
	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Router not found"})
	})
	router.Use(middleware.AuthMiddleware())
	router.POST("/expenses", h.CreateExpense)
	router.GET("/expenses/:id", h.GetExpenseById)
	router.PUT("/expenses/:id", h.UpdateExpenseById)
	router.GET("/expenses", h.GetAllExpenses)

	return router
}
