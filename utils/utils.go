package utils

import (
	"os"
	"strconv"
)

func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}