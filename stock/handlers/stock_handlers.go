package handlers

import "net/http"

type StockHandlers struct{}

func (s *StockHandlers) HandleGetStocks(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *StockHandlers) HandleListStocks(w http.ResponseWriter, r *http.Request) error {
	return nil
}
