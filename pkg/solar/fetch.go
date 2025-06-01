package solar

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PrivatePuffin/shem/pkg/helper"
)

// PanelConfig defines one set of solar panels' parameters
type PanelConfig struct {
	Latitude     float64
	Longitude    float64
	Declination  float64
	Azimuth      float64
	KilowattPeak float64
}

// Outputs accessible to other packages
var (
	SolarInput       []float64 // 24 hourly values, Wh per hour
	TotalSolarDay    float64   // sum of SolarInput Wh
	DailyUseSolarCov float64   // percentage of DailyUseEst covered by SolarInput
)

// Estimated daily use in Wh - set by the caller before FetchForecast
var DailyUseEst float64

const baseAPIURL = "https://api.forecast.solar/estimate"

type forecastSolarResponse struct {
	Result struct {
		WattHours struct {
			Today map[string]float64 `json:"today"`
		} `json:"watt_hours"`
	} `json:"result"`
}

// FetchForecast fetches solar forecasts for multiple panel configs, sums, and calculates coverage.
func FetchForecast(panels []PanelConfig) error {
	SolarInput = make([]float64, 24)

	for _, panel := range panels {
		url := fmt.Sprintf("%s/%.4f/%.4f/%.1f/%.1f/%.2f",
			baseAPIURL,
			panel.Latitude,
			panel.Longitude,
			panel.Declination,
			panel.Azimuth,
			panel.KilowattPeak,
		)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error fetching forecast for %+v: %w", panel, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("API returned status %s for %+v", resp.Status, panel)
		}

		var parsed forecastSolarResponse
		if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
			return fmt.Errorf("decoding error for %+v: %w", panel, err)
		}

		for hourStr, wh := range parsed.Result.WattHours.Today {
			t, err := time.Parse("15:04", hourStr)
			if err != nil {
				continue
			}
			helper.SolarInput[t.Hour()] += wh
		}
	}

	// Calculate total solar input for the day (Wh)
	helper.TotalSolarDay = 0
	for _, val := range SolarInput {
		helper.TotalSolarDay += val
	}

	// Calculate coverage percentage if DailyUseEst > 0
	if helper.DailyUseEst > 0 {
		helper.DailyUseSolarCov = (helper.TotalSolarDay / helper.DailyUseEst) * 100
	} else {
		helper.DailyUseSolarCov = 0
	}

	return nil
}

func Render() {
	panels := []PanelConfig{
		{Latitude: 52.52, Longitude: 13.405, Declination: 30, Azimuth: 180, KilowattPeak: 4.0},
		{Latitude: 52.52, Longitude: 13.405, Declination: 20, Azimuth: 90, KilowattPeak: 1.5},
	}

	if err := FetchForecast(panels); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hourly Solar Input (Wh): %v\n", helper.SolarInput)
	fmt.Printf("Total Solar Input Today (Wh): %.0f\n", helper.TotalSolarDay)
	fmt.Printf("Percentage of Estimated Daily Use Covered: %.2f%%\n", helper.DailyUseSolarCov)
}
