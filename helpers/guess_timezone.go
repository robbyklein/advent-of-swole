package helpers

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/initializers"
)

func GuessTimezone(r *http.Request) string {
	// Extract the client's IP address
	ip := getClientIP(r)
	if ip == "" {
		log.Println("Unable to extract client IP, defaulting to UTC")
		return "UTC"
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Printf("Invalid IP address: %v", ip)
		return "UTC"
	}

	// Perform the GeoIP lookup
	record, err := initializers.GeoDB.City(parsedIP)
	if err != nil {
		log.Printf("Error looking up IP in GeoLite2 database: %v", err)
		return "UTC"
	}

	// Return the timezone if available
	if record.Location.TimeZone != "" {
		return record.Location.TimeZone
	}

	// Default to UTC if no timezone is found
	return "UTC"
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Use the first IP in the list (real client IP)
		return strings.Split(forwarded, ",")[0]
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error parsing IP from RemoteAddr: %v", err)
		return ""
	}

	return ip
}

func geolocateTimezone(ip string) string {
	// String to ip
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Printf("Invalid IP address: %v", ip)
		return ""
	}

	record, err := initializers.GeoDB.City(parsedIP)
	if err != nil || record.Location.TimeZone == "" {
		log.Printf("Error looking up IP or missing timezone for IP: %v", ip)
		return ""
	}

	return record.Location.TimeZone
}

func isValidTimezone(tz string) bool {
	for _, valid := range config.Timezones {
		if strings.EqualFold(tz, valid) {
			return true
		}
	}
	return false
}
