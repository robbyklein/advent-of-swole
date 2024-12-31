package initializers

import (
	"fmt"

	"github.com/oschwald/geoip2-golang"
)

var GeoDB *geoip2.Reader

func LoadIPDatabase() {
	// Open the GeoLite2-City database
	var err error
	GeoDB, err = geoip2.Open("./GeoLite2-City.mmdb")

	if err != nil {
		fmt.Println("Failed to load geolocation database")
	}
}
