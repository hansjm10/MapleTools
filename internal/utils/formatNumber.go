// utils/formatNumber.go
package utils

import (
	"fmt"
	"math"
)

func FormatNumber(n float64) string {
	switch {
	case n >= 1e12:
		return fmt.Sprintf("%.2fT", n/1e12)
	case n >= 1e9:
		return fmt.Sprintf("%.2fB", n/1e9)
	case n >= 1e6:
		return fmt.Sprintf("%.2fM", n/1e6)
	case n >= 1e3:
		return fmt.Sprintf("%.2fK", n/1e3)
	default:
		return fmt.Sprintf("%.0f", math.Round(n))
	}
}
