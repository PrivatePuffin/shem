package solar

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type PanelConfig struct {
	Latitude     float64
	Longitude    float64
	Declination  float64
	Azimuth      float64
	KilowattPeak float64
}

var Forecast []float64 // Summed hourly forecast in Wh (length 24)

const baseAPIURL = "https://api.forecast.solar/estimate"

type forecastSolarResponse struct {
	Result struct {
		WattHours struct {
			Today map[string]float64 `json:"today"`
		} `json:"watt_hours"`
	} `json:"result"`
}

// FetchForecast fetches the solar forecast for multiple panel sets and sums the result
func FetchForecast(panels []PanelConfig) error {
	// Zero the forecast slice
	Forecast = make([]float64, 24)

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
			Forecast[t.Hour()] += wh
		}
	}

	return nil
}

func Render() {
	panels := []PanelConfig{
		{Latitude: 52.52, Longitude: 13.405, Declination: 30, Azimuth: 180, KilowattPeak: 4.0}, // south-facing
		{Latitude: 52.52, Longitude: 13.405, Declination: 20, Azimuth: 90, KilowattPeak: 1.5},  // east-facing
	}

	if err := FetchForecast(panels); err != nil {
		log.Fatal(err)
	}

	for i, val := range Forecast {
		fmt.Printf("Hour %02d: %.0f Wh\n", i, val)
	}
}
