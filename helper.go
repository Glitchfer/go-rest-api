package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/jsonapi"
)

func connect() *sql.DB {
	user := "root"
	password := ""
	host := "127.0.0.1"
	port := "3306"
	database := "konta_san"

	// root:secret@tcp(127.0.0.1:3306)/konta_san
	connection := fmt.Sprintf( /*format*/ "%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	db, err := sql.Open( /*driverName:*/ "mysql", connection)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func renderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set( /*key:*/ "Content-Type" /*value:*/, "application/json")
	jsonapi.MarshalPayload(w, data)
	// json.NewEncoder(w).Encode(data)
}
