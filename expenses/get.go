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

func (h *handler) GetAllExpenses(c *gin.Context) {
	var expenses []Expense
	rows, err := h.DB.Query("SELECT * FROM expenses")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.Title, &expense.Amount, &expense.Note, pq.Array(&expense.Tags))
		if err != nil {
			c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
			return
		}
		expenses = append(expenses, expense)
	}
	c.JSON(http.StatusOK, expenses)
}
