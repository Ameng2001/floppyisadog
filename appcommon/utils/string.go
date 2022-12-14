package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ttacon/libphonenumber"
)

// StringInSlice is a function similar to "x in y" Python construct
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// String to int64
func StringToInt(s string) int {
	int64Val, _ := strconv.ParseInt(s, 10, 64)
	return int(int64Val)
}

func StringToBool(s string) bool {
	boolVal, _ := strconv.ParseBool(s)
	return boolVal
}

// SpaceDelimitedStringNotGreater returns true if the first string
// is the same as the second string or does not contain any substring
// not contained in the second string (when split by space)
func SpaceDelimitedStringNotGreater(first, second string) bool {
	// Empty string is never greater
	if first == "" {
		return true
	}

	// Split the second string by space
	secondParts := strings.Split(second, " ")

	// Iterate over space delimited parts of the first string
	for _, firstPart := range strings.Split(first, " ") {
		// If the substring is not part of the second string, return false
		if !StringInSlice(firstPart, secondParts) {
			return false
		}
	}

	// The first string is the same or more restrictive
	// than the second string, return true
	return true
}

// Returns string encoded query string of the request
func GetQueryString(query url.Values) string {
	encoded := query.Encode()
	if len(encoded) > 0 {
		encoded = fmt.Sprintf("?%s", encoded)
	}
	return encoded
}

// formate phonenumber
const (
	defaultRegion = "CN" //
)

func ParseAndFormatPhonenumber(input string) (cleanPhonenumber string, err error) {
	// If empty string input - return empy string
	if input == "" {
		return
	}
	p, err := libphonenumber.Parse(input, defaultRegion)
	if err != nil {
		return "", fmt.Errorf("invalid phone number")
	}
	cleanPhonenumber = libphonenumber.Format(p, libphonenumber.E164)
	return
}
