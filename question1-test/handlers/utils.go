package handlers

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	otp := make([]byte, 4)
	for i := range otp {
		otp[i] = digits[rand.Intn(len(digits))]
	}
	return string(otp)
}

func IsOTPExpired(expirationTime time.Time) bool {
	// Define the layout to include only year, month, day, hour, minute, and second
	layout := "2006-01-02 15:04:05"

	// Format both expirationTime and the current time to this layout
	formattedExpirationTime := expirationTime.Format(layout)
	formattedCurrentTime := time.Now().Format(layout)

	// Parse them back into time.Time for comparison
	expirationTime, err := time.Parse(layout, formattedExpirationTime)
	if err != nil {
		fmt.Println("Error parsing formatted expiration time:", err)
		return false
	}

	currentTime, err := time.Parse(layout, formattedCurrentTime)
	if err != nil {
		fmt.Println("Error parsing formatted current time:", err)
		return false
	}

	fmt.Printf("Current time: %v, Expiration time: %v\n", currentTime, expirationTime)

	return currentTime.After(expirationTime)
}

func IsValidPhoneNumber(phone string) bool {
	// Regex to validate phone number format
	// Modify this regex to suit your specific requirements
	re := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}
