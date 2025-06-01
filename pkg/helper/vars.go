package helper

var (
	EntsoeAPIKey       string    = ""
	DailyUseEst        float64   = 30
	DailyUseSolarCov   float64   = 0
	SolarInput         []float64 = nil
	TotalSolarDay      float64   = 0
	PricesToday        []float64 = nil
	BuyPricesToday     []float64 = nil
	SellPricesToday    []float64 = nil
	TerugLevering      []float64 = nil
	ACCharge           []int     = nil
	MaxCharge          []int     = nil
	CheapestPrice      float64   = 0
	CheapestPriceIndex int       = 0
	Inkoopkosten       float64   = 0
	Heffingen          float64   = 0
	BTW                float64   = 21
	SellPriceMod       float64   = 0
	Logo                         = `
Shem

`
)
