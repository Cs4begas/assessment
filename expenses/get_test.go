//go:build unit
// +build unit

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

func TestGetAllExpenses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"})
	expected_expenses := []Expense{
		{ID: 1, Title: "Groceries", Amount: 100, Note: "Food for the week", Tags: []string{"food", "necessities"}},
		{ID: 2, Title: "Rent", Amount: 1000, Note: "Rent for the month", Tags: []string{"rent", "necessities"}},
	}
	for _, e := range expected_expenses {
		rows.AddRow(e.ID, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	}
	mock.ExpectQuery("SELECT \\* FROM expenses").WillReturnRows(rows)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)

	h.GetAllExpenses(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
	var expense []Expense
	json.Unmarshal(rr.Body.Bytes(), &expense)

	for i, e := range expected_expenses {
		assert.Equal(t, e.ID, expense[i].ID)
		assert.Equal(t, e.Title, expense[i].Title)
		assert.EqualValues(t, e.Amount, expense[i].Amount)
		assert.Equal(t, e.Note, expense[i].Note)
		assert.Equal(t, e.Tags, expense[i].Tags)
	}

}

func TestGetAllExpensesQueryError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	mock.ExpectQuery("SELECT \\* FROM expenses").WillReturnError(sql.ErrConnDone)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)

	h.GetAllExpenses(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	var errResp Err
	json.Unmarshal(rr.Body.Bytes(), &errResp)
	assert.NotNil(t, errResp.Message)
}

func TestGetAllExpensesScanError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	rows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(nil, "Groceries", 100, "Food for the week", pq.Array([]string{"food", "necessities"})).
		AddRow(1, "Rent", 1000, "Rent for the month", pq.Array([]string{"rent", "necessities"}))
	mock.ExpectQuery("SELECT \\* FROM expenses").WillReturnRows(rows)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)

	h.GetAllExpenses(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	var errResp Err
	json.Unmarshal(rr.Body.Bytes(), &errResp)
	assert.NotNil(t, errResp.Message)
}
