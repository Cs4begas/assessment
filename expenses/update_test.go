package expenses

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	mock.ExpectExec("UPDATE expenses SET").WillReturnResult(sqlmock.NewResult(1, 1))

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	jsonReq := `{"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`
	c.Request, _ = http.NewRequest("PUT", "/expenses/1", bytes.NewBufferString(jsonReq))

	h.UpdateExpenseById(c)

	assert.Equal(t, http.StatusOK, c.Writer.Status())
	var expense Expense
	json.Unmarshal(rr.Body.Bytes(), &expense)
	assert.Equal(t, "apple smoothie", expense.Title)
	assert.EqualValues(t, 89, expense.Amount)
	assert.Equal(t, "no discount", expense.Note)
	assert.Equal(t, []string{"beverage"}, expense.Tags)

}

func TestUpdateExpenseByIdFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	mock.ExpectExec("UPDATE expenses SET").WillReturnError(sql.ErrNoRows)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	jsonReq := `{"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`
	c.Request, _ = http.NewRequest("PUT", "/expenses/1", bytes.NewBufferString(jsonReq))

	h.UpdateExpenseById(c)

	assert.Equal(t, http.StatusInternalServerError, c.Writer.Status())
	var errResp Err
	json.Unmarshal(rr.Body.Bytes(), &errResp)
	assert.NotNil(t, errResp.Message)
}

func TestUpdateExpenseInvalidJson(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	h := NewApplication(db)

	rr := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rr)
	jsonReq := `{"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]`
	c.Request, _ = http.NewRequest("PUT", "/expenses/1", bytes.NewBufferString(jsonReq))

	h.UpdateExpenseById(c)

	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
	var errResp Err
	json.Unmarshal(rr.Body.Bytes(), &errResp)
	assert.NotNil(t, errResp.Message)
}
