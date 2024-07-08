package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"log"
	"io"
	"errors"
	"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/cotacao", dollarQuotationHandler)
	http.ListenAndServe(":8080", nil)
}

func dollarQuotationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	resp, err := dollarQuotation(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(resp))
	//w.Write([]byte(`Hello world`))

	err = createDataBase()
	if err != nil {
		log.Output(1, err.Error())
		return
	}

}

func createDataBase() error {
	file, err := os.Create("./database.db")
	if err != nil {
		return err
	}
	file.Close()

	sqliteDatabase, err := sql.Open("sqlite3", "./database.db")
	defer sqliteDatabase.Close()
	if err != nil {
		return err
	}
	err = createTable(sqliteDatabase)
	if err != nil {
		return err
	}
	return nil
}

func createTable(db *sql.DB) error {
	createTableSQL := `create table quotations ("quotation" TEXT);`
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Output(1, err.Error())
		return err
	}
	statement.Exec()
	insertCode := `insert into quotations(quotation) values(?)`
	statement, err = db.Prepare(insertCode)
	if err != nil {
		return err
	}
	_, err = statement.Exec("xxx")
	if err != nil {
		return err
	}
	return nil
}



func dollarQuotation(ctx context.Context) (string, error) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	res, err := http.Get(url)
	if err != nil {
		log.Output(1, fmt.Sprintf("Error connecting to %s: %s\n", url, err))
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Output(1, err.Error())
		return "", err
	}
	res.Body.Close()
	select {
	case <-ctx.Done():
		err := errors.New("API call timed out")
		log.Output(1, err.Error())
		return "", err
	default:
		log.Output(1, string(body))
		return string(body), nil
	}
}
