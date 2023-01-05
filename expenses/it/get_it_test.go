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

func TestITGetExpenseById(t *testing.T) {
	//Arrange Post expense
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

	var expensePost expenses.Expense
	json.Unmarshal(byteBody, &expensePost)

	//Get expense
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, expensePost.ID), nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")

	resp, _ = callRequest(req)

	byteBody, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var expenseGet expenses.Expense
	json.Unmarshal(byteBody, &expenseGet)
	assert.Equal(t, expensePost.ID, expenseGet.ID)
	assert.Equal(t, expensePost.Title, expenseGet.Title)
	assert.Equal(t, expensePost.Amount, expenseGet.Amount)
	assert.Equal(t, expensePost.Note, expenseGet.Note)
	assert.Equal(t, expensePost.Tags, expenseGet.Tags)

	shutdownServer(srv, t)
}

func TestITGetExpenseByIdUnauthorized(t *testing.T) {
	srv := startServer()
	connectServer()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), nil)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")

	//ACT
	resp, _ := callRequest(req)

	//Assert
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	shutdownServer(srv, t)
}

func TestITGetAllExpenses(t *testing.T) {
	srv := startServer()
	connectServer()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/", serverPort), nil)
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")

	resp, _ := callRequest(req)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var expenses []expenses.Expense
	json.Unmarshal(byteBody, &expenses)
	assert.NotEqual(t, 0, len(expenses))

	shutdownServer(srv, t)
}

func TestITGetAllExpensesUnauthorized(t *testing.T) {
	srv := startServer()
	connectServer()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/1", serverPort), nil)
	assert.NoError(t, err)
	req.Header.Add("Content-Type", "application/json")

	//ACT
	resp, _ := callRequest(req)

	//Assert
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	shutdownServer(srv, t)
}
