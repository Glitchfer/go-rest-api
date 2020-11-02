package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/jsonapi"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    int64  `jsonapi:"primary,products"`
	Name  string `jsonapi:"attr,product_name"`
	Price int64  `jsonapi:"attr,product_price"`
}

var mysqlDB *sql.DB

func Server() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleHome).Methods("GET")
	router.HandleFunc("/api/products", BrowseProduct).Methods("GET")
	router.HandleFunc("/api/products", CreateProduct).Methods("POST")
	return router

}

func main() {
	mysqlDB = connect()
	defer mysqlDB.Close()

	router := Server()

	log.Println(`API sudah berjalan di port:8080`)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleHome(writer http.ResponseWriter, response *http.Request) {
	writer.Header().Set( /*key:*/ "Content-type" /*value:*/, "application/json")
	writer.WriteHeader(http.StatusOK)

	json.NewEncoder(writer).Encode(map[string]interface{}{
		"status": 200,
		"msg":    "Hello world",
	})
}

func BrowseProduct(writer http.ResponseWriter, request *http.Request) {
	rows, err := mysqlDB.Query("SELECT product_id, product_name, product_price FROM product")
	if err != nil {
		renderJSON(writer, map[string]interface{}{
			"msg": "Product not found",
		})
	}

	var products []*Product

	for rows.Next() {
		var product Product

		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Print(err)
		} else {
			products = append(products, &product)
		}
	}

	renderJSON(writer, products)
}

func CreateProduct(writer http.ResponseWriter, request *http.Request) {
	var product Product

	err := jsonapi.UnmarshalPayload(request.Body, &product)
	if err != nil {
		log.Print(err)
		return
	}

	query, err := mysqlDB.Prepare("INSERT INTO product (product_name, product_price) VALUES (?, ?)")
	if err != nil {
		log.Print(err)
		return
	}

	result, err := query.Exec(product.Name, product.Price)
	if err != nil {
		log.Print(err)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return
	}

	product.ID = lastID
	writer.WriteHeader(http.StatusCreated)
	jsonapi.MarshalPayload(writer, &product)
}
