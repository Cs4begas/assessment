package expenses

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseById(c *gin.Context) {
	id := c.Param("id")
	var expense Expense
	row := h.DB.QueryRow("SELECT * FROM expenses WHERE id=$1", id)
	err := row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, Err{Message: err.Error()})
	}
	c.JSON(http.StatusOK, expense)
}
