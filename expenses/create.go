package expenses

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func CreateExpense(c *gin.Context) {
	var expense Expense

	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	db := GetDB()
	sqlStatement := `INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id`
	row := db.QueryRow(sqlStatement, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))
	err := row.Scan(&expense.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	fmt.Printf("id : % #v\n", expense)

	c.JSON(http.StatusCreated, expense)

}
