package power

import "github.com/PrivatePuffin/shem/pkg/helper"

// SetACCharge processes helper.BuyPricesToday:
// if the price is negative, set to 100, else set to 0
func SetACCharge() {
	for i, price := range helper.BuyPricesToday {
		if price < 0 {
			helper.ACCharge[i] = 100
		} else {
			helper.ACCharge[i] = 0
		}
	}
}
