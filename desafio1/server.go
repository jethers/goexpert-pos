package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

func main() {
	myServer := http.NewServeMux()
	myServer.HandleFunc("/cotacao", dollarQuotationHandler)
	http.ListenAndServe(":8080", myServer)
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

	ctx, cancel = context.WithTimeout(ctx, time.Millisecond*10)
	err = dollarQuotationDataBase(ctx, resp)
	if err != nil {
		log.Output(1, err.Error())
		return
	}
}

func dollarQuotationDataBase(ctx, resp string) (error){
    db, err := sql.Open("sqlite3", "./quotation.db")
	if err != nil{
		log.Output(1, err.Error())
	}
	defer db.Close()

	select {
	case <-ctx.Done():
		err := errors.New("Database call timed out")
		log.Output(1, err.Error())
		return "", err
	default:
      return
	}
}

func dollarQuotation(ctx context.Context) (string, error) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	res, err := http.Get(url)
	if err != nil || res.StatusCode != http.StatusOK {
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
