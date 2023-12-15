package util

const (
	CNY = "CNY"
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CNY, USD, EUR, CAD:
		return true
	}
	return false
}
