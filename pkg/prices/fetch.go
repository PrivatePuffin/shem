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

var pricesToday, pricesTomorrow []float64

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
func Fetch() {
	now := time.Now().UTC()
	startToday := now.Format("200601020000")
	endTomorrow := now.Add(48 * time.Hour).Format("200601020000")

	allPrices, err := fetchPrices(startToday, endTomorrow)
	if err != nil {
		log.Error().Msgf("Error fetching prices: %v", err)
		os.Exit(1)
	}

	// split into today (first 24) and tomorrow (next 24) if available
	if len(allPrices) >= 24 {
		helper.PricesToday = allPrices[:24]
	}

	if len(allPrices) >= 48 {
		helper.PricesTomorrow = allPrices[24:48]
		helper.PricesTomorrowAvailable = true
	} else {
		helper.PricesTomorrow = nil
		helper.PricesTomorrowAvailable = false
	}

	// For demonstration only:
	fmt.Println("Prices today:", helper.PricesToday)
	if helper.PricesTomorrowAvailable {
		fmt.Println("Prices tomorrow:", helper.PricesTomorrow)
		fmt.Println("Tomorrow available?", helper.PricesTomorrowAvailable)
	}
}
