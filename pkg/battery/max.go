package battery

import "github.com/PrivatePuffin/shem/pkg/helper"

// AdjustMaxCharge updates helper.MaxCharge with logic:
// - Default 100 for all indices
// - For indices 7–14:
//   - Apply sell price zero checks to reduce MaxCharge
//
// - Then override: if ACCharge[i] != 0, force MaxCharge[i] = 100
func AdjustMaxCharge() {
	n := len(helper.SellPricesToday)
	if len(helper.MaxCharge) != n {
		helper.MaxCharge = make([]int, n)
	}
	if len(helper.ACCharge) != n {
		helper.ACCharge = make([]int, n)
	}

	// Step 1: Default all to 100
	for i := range helper.MaxCharge {
		helper.MaxCharge[i] = 100
	}

	// Step 2: Apply zero-price logic to indices 7–14 only
	for i := 7; i <= 14 && i < n; i++ {
		// Apply logic only if ACCharge is zero (will be overridden later if not)
		if helper.SellPricesToday[i] == 0 {
			helper.MaxCharge[i] = 20
		}
		if i+1 < n && helper.SellPricesToday[i+1] == 0 {
			helper.MaxCharge[i] = 20
		}
		if i+2 < n && helper.SellPricesToday[i+2] == 0 {
			helper.MaxCharge[i] = 30
		}
		if i+3 < n && helper.SellPricesToday[i+3] == 0 {
			helper.MaxCharge[i] = 40
		}
		if i+4 < n && helper.SellPricesToday[i+4] == 0 {
			helper.MaxCharge[i] = 50
		}
		if i+5 < n && helper.SellPricesToday[i+5] == 0 {
			helper.MaxCharge[i] = 60
		}
	}

	// Step 3: Override with ACCharge: if != 0 → force to 100
	for i := 0; i < n; i++ {
		if helper.ACCharge[i] != 0 {
			helper.MaxCharge[i] = 100
		}
	}
}
