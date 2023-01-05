//go:build integration
// +build integration

package expenses_it_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/cs4begas/assessment/expenses"
	"github.com/stretchr/testify/assert"
)

func TestITCreateExpense(t *testing.T) {
	srv := startServer()
	connectServer()

	reqBody := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath",
			"tags": ["food", "beverage"]
		}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")

	resp, _ := callRequest(req)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var expense expenses.Expense
	json.Unmarshal(byteBody, &expense)
	assert.Equal(t, "strawberry smoothie", expense.Title)
	assert.Equal(t, 79.0, expense.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", expense.Note)
	assert.Equal(t, []string{"food", "beverage"}, expense.Tags)

	shutdownServer(srv, t)

}

func TestITCreateExpenseUnauthorized(t *testing.T) {
	//Arrange
	srv := startServer()
	connectServer()

	reqBody := `{
			"title": "strawberry smoothie",
			"amount": 79,
			"note": "night market promotion discount 10 bath",
			"tags": ["food", "beverage"]
		}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")

	//ACT
	resp, _ := callRequest(req)

	//Assert
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	shutdownServer(srv, t)
}
