package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/entity"
	"server/infra"
	"server/service"
	"time"
)

func main() {
	fmt.Println("Server running")
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetExchangeHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func GetExchangeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()
	exchangeRate, err := service.GetExchangeRate(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	data, err := PersistData(exchangeRate)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

}

func PersistData(ex *service.ExchangeRate) (entity.Bid, error) {
	db, err := infra.OpenDatabase()
	if err != nil {
		return entity.Bid{}, nil
	}
	defer db.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	service := service.Data{Bid: entity.Bid{BidValue: ex.USDBRL.Bid}}
	err = service.InsertBID(ctx, db)
	if err != nil {
		return entity.Bid{}, err
	}
	return service.Bid, nil
}
