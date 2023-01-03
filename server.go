package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	fmt.Println("start at port:", os_port)
	log.Fatal(router.Run(":" + os_port))
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
