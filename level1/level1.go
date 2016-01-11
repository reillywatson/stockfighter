package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const apiKey = "6506b277dfc546403ff7df2c3577d5af48e424b0"

type Direction string

const (
	Buy  Direction = "buy"
	Sell Direction = "sell"
)

type OrderType string

const (
	Limit  OrderType = "limit"
	Market OrderType = "market"
)

type BuySellRequest struct {
	Account   string    `json:"account"`
	Venue     string    `json:"venue"`
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price,omitempty"`
	Quantity  int       `json:"qty"`
	Direction Direction `json:"direction"`
	OrderType OrderType `json:"orderType"`
}

func main() {
	err := buyStock()
	if err != nil {
		fmt.Println(err)
	}
}

func buyStock() error {
	order := BuySellRequest{
		Account:   "LOB39707345",
		Venue:     "EWBEX",
		Symbol:    "LSIM",
		Quantity:  100,
		Direction: Buy,
		OrderType: Market,
	}
	body, err := json.Marshal(order)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks/%s/orders", order.Venue, order.Symbol), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Add("X-Starfighter-Authorization", apiKey)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respBody))
	return nil
}
