// unit test
package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_Home(t *testing.T) {
	request, _ := http.NewRequest( /*method:*/ "GET" /*url:*/, "/" /*body:*/, nil)
	response := httptest.NewRecorder()

	Server().ServeHTTP(response, request)
	expectedResponse := `{"msg":"Hello world","status":200}`
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Log(err)
	}

	// Assert status = cek status benar / tidak
	assert.Equal(t /*expected:*/, 200, response.Code /*msgAndArgs...:*/, "Invalid response code")

	// Assert response
	assert.Equal(t, expectedResponse, string(bytes.TrimSpace(responseBody)))
}

func Test_BrowseProduct(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Log(err)
	}
	rows := sqlmock.NewRows([]string{"product_id", "product_name", "product_price"}).
		AddRow(1, "Espresso", 10000)
	mock.ExpectQuery("SELECT product_id, product_name, product_price FROM product").
		WillReturnRows(rows)

	mysqlDB = db
	request, _ := http.NewRequest( /*method:*/ "GET" /*url:*/, "/api/products" /*body:*/, nil)
	response := httptest.NewRecorder()

	Server().ServeHTTP(response, request)
	expectedResponse := `{"data":[{"type":"products","id":"1","attributes":{"product_name":"Espresso","product_price":10000}}]}`
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Log(err)
	}

	assert.Equal(t /*expected:*/, 200, response.Code /*msgAndArgs...:*/, "Invalid response code")
	assert.Equal(t, expectedResponse, string(bytes.TrimSpace(responseBody)))
}

func Test_CrateProduct(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Log(err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO product (product_name, product_price) VALUES (?, ?)")
	mock.ExpectExec("INSERT INTO product (product_name, product_price) VALUES (?, ?)").
		WithArgs("Lemon Tea", 7500).
		WillReturnResult(sqlmock.NewResult(87, 1))

	mysqlDB = db

	data := map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": map[string]interface{}{
				"product_name":  "Lemon Tea",
				"product_price": 7500,
			},
		},
	}

	requestBody, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	Server().ServeHTTP(response, request)
	expectedResponse := `{"data":{"type":"products","id":"87","attributes":{"product_name":"Lemon Tea","product_price":7500}}}`
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Log(err)
	}
	t.Log(response.Body.String())
	assert.Equal(t, http.StatusCreated, response.Code /*msgAndArgs...:*/, "Invalid response code")
	assert.Equal(t, expectedResponse, string(bytes.TrimSpace(responseBody)))
}
