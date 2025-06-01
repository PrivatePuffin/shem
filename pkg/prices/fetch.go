package prices

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/PrivatePuffin/shem/pkg/helper"
	"github.com/rs/zerolog/log"
)

var (
	securityToken = helper.EntsoeAPIKey
)

const (
	apiURL       = "https://web-api.tp.entsoe.eu/api"
	nlDomain     = "10YNL----------L"
	documentType = "A44"
)

type TimeSeries struct {
	Period struct {
		Points []struct {
			Position int     `xml:"position"`
			Price    float64 `xml:"price.amount"`
		} `xml:"Point"`
	} `xml:"period"`
}

type Response struct {
	TimeSeries []TimeSeries `xml:"TimeSeries"`
}

func fetchPrices(start, end string) ([]float64, error) {
	if securityToken == "" {
		return nil, errors.New("ENTSOE API token is missing")
	}
	reqURL := fmt.Sprintf(
		"%s?securityToken=%s&documentType=%s&in_Domain=%s&out_Domain=%s&periodStart=%s&periodEnd=%s",
		apiURL, securityToken, documentType, nlDomain, nlDomain, start, end)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s\n%s", resp.Status, string(body))
	}

	var data Response
	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var prices []float64
	for _, ts := range data.TimeSeries {
		for _, p := range ts.Period.Points {
			prices = append(prices, p.Price/1000) // €/kWh
		}
	}
	return prices, nil
}

// BuyPriceToday calculates the corrected prices for today including fees and tax (BTW)
func BuyPriceToday(pricesToday []float64) []float64 {
	taxFactor := 1 + helper.BTW/100.0 // ensure float division
	corrected := make([]float64, len(pricesToday))
	for i, p := range pricesToday {
		sum := p + helper.Inkoopkosten + helper.Heffingen
		if sum < 0 {
			corrected[i] = sum
		} else {
			corrected[i] = sum * taxFactor
		}
	}
	return corrected
}
func SellPriceToday(pricesToday []float64) []float64 {
	sellPrices := make([]float64, len(pricesToday))
	for i, p := range pricesToday {
		sp := p + helper.SellPriceMod
		if sp < 0 {
			sp = 0
		}
		sellPrices[i] = sp
	}
	return sellPrices
}

// FindCheapestPrice returns the lowest price and its index in the slice
func FindCheapestPrice(prices []float64) (float64, int) {
	if len(prices) == 0 {
		return 0, -1 // no prices available
	}
	minPrice := prices[0]
	minIndex := 0
	for i, p := range prices {
		if p < minPrice {
			minPrice = p
			minIndex = i
		}
	}
	return minPrice, minIndex
}

func Fetch() {
	now := time.Now().UTC()
	startToday := now.Truncate(24 * time.Hour).Format("200601020000")
	endToday := now.Truncate(24 * time.Hour).Add(24 * time.Hour).Format("200601020000")

	pricesToday, err := fetchPrices(startToday, endToday)
	if err != nil {
		log.Error().Msgf("Error fetching today's prices: %v", err)
		os.Exit(1)
	}
	helper.PricesToday = pricesToday

	helper.BuyPricesToday = BuyPriceToday(pricesToday)
	helper.SellPricesToday = SellPriceToday(pricesToday)

	price, idx := FindCheapestPrice(pricesToday)
	helper.CheapestPrice = price
	helper.CheapestPriceIndex = idx
}

func Render() {
	fmt.Println("Prices today:", helper.PricesToday)
	fmt.Println("Buy Prices today:", helper.BuyPricesToday)
	fmt.Println("Sell Prices today:", helper.SellPricesToday)

	if helper.CheapestPriceIndex >= 0 {
		fmt.Printf("Cheapest price: %.4f €/kWh at index %d\n", helper.CheapestPrice, helper.CheapestPriceIndex)
	} else {
		fmt.Println("No prices available to find the cheapest price.")
	}
}
