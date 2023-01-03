package expenses

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetExpenseById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(1, "Groceries", 100, "Food for the week", pq.Array([]string{"food", "necessities"}))

	mock.ExpectQuery("SELECT (.+) FROM expenses").WithArgs("1").WillReturnRows(rows)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	h.GetExpenseById(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
	var expense Expense
	json.Unmarshal(rr.Body.Bytes(), &expense)
	assert.Equal(t, 1, expense.ID)
	assert.Equal(t, "Groceries", expense.Title)
	assert.EqualValues(t, 100, expense.Amount)
	assert.Equal(t, "Food for the week", expense.Note)
	assert.Equal(t, []string{"food", "necessities"}, expense.Tags)

}

func TestGetExpenseByIdNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	mock.ExpectQuery("SELECT (.+) FROM expenses").WithArgs("1").WillReturnError(sql.ErrNoRows)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	h.GetExpenseById(c)

	assert.Equal(t, http.StatusNotFound, c.Writer.Status())
	var errResp Err
	json.Unmarshal(rr.Body.Bytes(), &errResp)
	assert.Equal(t, "Not found expense", errResp.Message)

}

func TestGetExpenseByIdError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	mock.ExpectQuery("SELECT (.+) FROM expenses").WithArgs("1").WillReturnError(sql.ErrConnDone)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	h.GetExpenseById(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	var errResp Err
	json.Unmarshal(rr.Body.Bytes(), &errResp)
	assert.NotNil(t, errResp.Message)
}
