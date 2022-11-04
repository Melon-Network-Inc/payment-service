package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)


const float64EqualityThreshold = 1e-9
const float64ShortenMinThreshold = 1e-5
const float64ShortenMaxThreshold = 1e6

func AlmostEqual(a, b float64) bool {
    return math.Abs(a - b) <= float64EqualityThreshold
}

func GetFloatPointString(floatPoint float64) string {
	if AlmostEqual(floatPoint, 0) {
		return "0"
	}
	if AlmostEqual(floatPoint, float64ShortenMinThreshold) {
		return strconv.FormatFloat(floatPoint, 'e', -1, 5)
	}
	if floatPoint > float64ShortenMaxThreshold {
		return strconv.FormatFloat(floatPoint, 'e', -1, 5)
	}
	return fmt.Sprintf("%f", floatPoint)
}

func Shorten(number string) string {
	if strings.Contains(number, "e+00") {
		return number[0:len(number) - 4]
	}
	return number
}