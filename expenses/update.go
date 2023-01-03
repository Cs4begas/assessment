package expenses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *Handler) UpdateExpenseById(c *gin.Context) {
	var expense Expense
	id := c.Param("id")
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sqlStatement := `UPDATE expenses SET title=$1, amount=$2, note=$3, tags=$4 WHERE id=$5`
	_, err := h.DB.Exec(sqlStatement, expense.Title, expense.Amount, expense.Note, pq.Array(expense.Tags), id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Err{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, expense)
}
