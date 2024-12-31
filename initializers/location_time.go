package initializers

import (
	"fmt"
	"time"
)

var Location *time.Location

func LoadLocation() {
	var err error

	Location, err = time.LoadLocation("America/Los_Angeles")

	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
}
