package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const apiURL = "https://api.exchangerate-api.com/v4/latest/USD" 

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

func fetchRates() (ExchangeRates, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return ExchangeRates{}, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	var rates ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return ExchangeRates{}, fmt.Errorf("failed to decode exchange rates: %w", err)
	}

	return rates, nil
}

func convert(amount float64, from string, to string, rates ExchangeRates) (float64, error) {
	fromRate, ok := rates.Rates[from]
	if !ok {
		return 0, fmt.Errorf("invalid currency code: %s", from)
	}

	toRate, ok := rates.Rates[to]
	if !ok {
		return 0, fmt.Errorf("invalid currency code: %s", to)
	}

	return amount / fromRate * toRate, nil
}

func parseArguments() (float64, string, string, error) {
	if len(os.Args) != 4 {
		return 0, "", "", fmt.Errorf("usage: go run main.go <amount> <from_currency> <to_currency>")
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("invalid amount: %s", os.Args[1])
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	return amount, fromCurrency, toCurrency, nil
}

func main() {
	amount, fromCurrency, toCurrency, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	rates, err := fetchRates()
	if err != nil {
		fmt.Printf("Error fetching rates: %v\n", err)
		return
	}

	result, err := convert(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Error converting currency: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}
