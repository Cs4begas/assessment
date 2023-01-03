package expenses

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h *handler) GetExpenseById(c *gin.Context) {
	var expense Expense
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, Err{Message: "param id is required"})
		return
	}
	row := h.DB.QueryRow("SELECT * FROM expenses WHERE id=$1", id)
	err := row.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, Err{Message: "Not found expense"})
			return
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, Err{Message: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, expense)
}
