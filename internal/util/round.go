package util

import "math"

func Round(value, roundOn, decimalDigits float64) float64 {

	var round float64
	pow := math.Pow(10, decimalDigits)
	digit := pow * value
	_, div := math.Modf(digit)

	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}

	return round / pow
}
