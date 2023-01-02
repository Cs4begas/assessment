package expenses_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cs4begas/assessment/expenses"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	h := expenses.NewApplication(db)
	router := gin.Default()
	router.POST("/expenses", h.CreateExpense)

	rs1 := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("INSERT INTO expenses").WithArgs("grocery", 20.5, "food", pq.Array([]string{"food", "grocery"})).WillReturnRows(rs1)

	tests := []struct {
		name string
		json string
		code int
	}{
		{
			name: "invalid json",
			json: `{"title":"grocery","amount":20.5,"note":"food","tags":["food","grocery"]`,
			code: http.StatusBadRequest,
		},
		{
			name: "valid json",
			json: `{"title":"grocery","amount":20.5,"note":"food","tags":["food","grocery"]}`,
			code: http.StatusCreated,
		},
		{
			name: "http internal server error",
			json: `{"title":"httpInternalServerError","amount":20.5,"note":"food","tags":["food","grocery"]}`,
			code: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/expenses", bytes.NewBufferString(test.json))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, test.code, rr.Code)
		})
	}
}
