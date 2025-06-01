package power

import (
	"fmt"

	"github.com/PrivatePuffin/shem/pkg/helper"
)

// Teruglevering calculates for each price in pricesToday a value between 0 and 100
// If the price is > 0, teruglevering is 100
// If the price is 0, teruglevering is 0
func Teruglevering(pricesToday []float64) []float64 {
	result := make([]float64, len(pricesToday))
	for i, p := range pricesToday {
		if p > 0 {
			result[i] = 100
		} else if p == 0 {
			result[i] = 0
		}
		// if price < 0, behavior not defined, you can adjust if needed
	}
	return result
}

// Render prints the teruglevering slice to the console
func Fetch() {
	terug := Teruglevering(helper.PricesToday)
	fmt.Println("Teruglevering values:", terug)

}
