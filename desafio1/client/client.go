package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()
	resp, err := callQuotation(ctx)
	if err != nil {
	    panic(err)
	}

    result := make(map[string]interface{})
	json.Unmarshal([]byte(resp), &result)
	dollar := result["USDBRL"].(map[string]interface{})
	cotacao := fmt.Sprintf("Dólar:{%v}", dollar["bid"])
   
	err = os.WriteFile("./data/cotacao.txt", []byte(cotacao), 0644)
	if err != nil {
	    panic(err)
	}
}

func callQuotation(ctx context.Context) (string, error) {
	url := "http://localhost:8080/cotacao"
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
		err := errors.New("Server API call timed out")
		log.Output(1, err.Error())
		return "", err
	default:
		log.Output(1, string(body))
		return string(body), nil
	}
}
