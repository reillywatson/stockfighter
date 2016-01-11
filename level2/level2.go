package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
	Price     int       `json:"price,omitempty"`
	Quantity  int       `json:"qty"`
	Direction Direction `json:"direction"`
	OrderType OrderType `json:"orderType"`
}

type QuoteRequest struct {
	Account string `json:"account"`
	Venue   string `json:"venue"`
	Symbol  string `json:"symbol"`
}

func main() {
	for {
		account := "FFS11174186"
		venue := "HKNBEX"
		symbol := "AWH"
		target := 3595
		bid, ask, err := getQuote(account, venue, symbol)
		fmt.Println("BID:", bid, "ASK:", ask)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if ask < target {
			err = buyStock(account, venue, symbol, ask+1, 1000, Limit)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		time.Sleep(time.Second)
	}
}

var client = http.Client{}

func getQuote(account, venue, symbol string) (int, int, error) {
	quote := QuoteRequest{
		Account: account,
		Venue:   venue,
		Symbol:  symbol,
	}
	body, err := json.Marshal(quote)
	if err != nil {
		return 0, 0, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.stockfighter.io/ob/api/venues/%s/stocks/%s", venue, symbol), bytes.NewBuffer(body))
	if err != nil {
		return 0, 0, err
	}
	req.Header.Add("X-Starfighter-Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	fmt.Println(string(respBody))
	var parsed struct {
		Bids []struct {
			Price    int `json:"price"`
			Quantity int `json:"qty"`
		} `json:"bids"`
		Asks []struct {
			Price    int `json:"price"`
			Quantity int `json:"qty"`
		} `json:"asks"`
	}
	json.Unmarshal(respBody, &parsed)
	lowest := 1000000
	highest := 0
	for _, ask := range parsed.Asks {
		if ask.Price < lowest {
			lowest = ask.Price
		}
	}
	for _, bid := range parsed.Bids {
		if bid.Price > highest {
			highest = bid.Price
		}
	}
	return highest, lowest, nil
}

func buyStock(account, venue, symbol string, price, quantity int, orderType OrderType) error {
	order := BuySellRequest{
		Account:   account,
		Venue:     venue,
		Symbol:    symbol,
		Price:     price,
		Quantity:  quantity,
		Direction: Buy,
		OrderType: orderType,
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
