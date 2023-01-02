package expenses

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *handler) CreateExpense(c *gin.Context) {
	expense := Expense{}

	if err := c.ShouldBindJSON(&expense); err != nil {
		log.Println("Hello")
		c.AbortWithStatusJSON(http.StatusBadRequest, Err{Message: err.Error()})
		return
	}

	db := h.DB
	sqlStatement := `INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err := db.Exec(sqlStatement, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)

}
