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

func TestUpdateExpenseById(t *testing.T) {
	//Arrange Post expense
	srv := startServer()
	connectServer()

	reqBodyPost := `{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBodyPost))
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

	reqBodyPut := `{
		"title": "milk tea",
		"amount": 1000,
		"note": "market promotion discount 10 bath",
		"tags": ["sushi", "beverage"]
	}`

	req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, expensePost.ID), strings.NewReader(reqBodyPut))
	assert.NoError(t, err)
	req.Header.Add("Authorization", "November 10, 2009")
	req.Header.Add("Content-Type", "application/json")

	resp, _ = callRequest(req)

	byteBody, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var expected_expenses expenses.Expense
	json.Unmarshal([]byte(reqBodyPut), &expected_expenses)

	var expensePut expenses.Expense
	json.Unmarshal(byteBody, &expensePut)

	fmt.Printf("expensesPut: %+v", expensePut)

	assert.Equal(t, expected_expenses.Title, expensePut.Title)
	assert.Equal(t, expected_expenses.Amount, expensePut.Amount)
	assert.Equal(t, expected_expenses.Note, expensePut.Note)
	assert.Equal(t, expected_expenses.Tags, expensePut.Tags)

	shutdownServer(srv, t)
}
