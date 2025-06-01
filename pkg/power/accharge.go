package power

import (
	"sort"

	"github.com/PrivatePuffin/shem/pkg/helper"
)

// SetACCharge processes helper.BuyPricesToday:
// if the price is negative, set to 100, else set to 0
func SetACCharge() {
	// Reset all ACCharge to 0 first
	for i := range helper.ACCharge {
		helper.ACCharge[i] = 0
	}

	// Collect negative prices with their indexes
	type priceIndex struct {
		Index int
		Price float64
	}
	var negatives []priceIndex
	for i, price := range helper.BuyPricesToday {
		if price < 0 {
			negatives = append(negatives, priceIndex{i, price})
		}
	}

	// Sort by price ascending (lowest first)
	sort.Slice(negatives, func(i, j int) bool {
		return negatives[i].Price < negatives[j].Price
	})

	// Set ACCharge to 100 for the lowest 2 negative prices
	for i := 0; i < len(negatives) && i < 2; i++ {
		helper.ACCharge[negatives[i].Index] = 100
	}
}
