package power

import (
	"fmt"

	"github.com/PrivatePuffin/shem/pkg/helper"
)

// Render prints the teruglevering slice to the console
func Gen() {
	Teruglevering(helper.PricesToday)
	Teruglevering(helper.PricesToday)

}

// Render prints the teruglevering slice to the console
func Render() {
	Gen()
	fmt.Println("Teruglevering:", helper.TerugLevering)
	fmt.Println("AC charge:", helper.ACCharge)

}
