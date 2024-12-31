package controllers

import (
	"net/http"
	"time"

	"github.com/robbyklein/swole/helpers"
	"github.com/robbyklein/swole/initializers"
)

func HomeGET(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().In(initializers.Location)
	// currentDay := currentTime.Day()
	currentMonth := int(currentTime.Month())
	currentYear := currentTime.Year()
	currentDay := 12
	annualInt := 2025 - currentYear
	annual := helpers.OrdinalWords(annualInt)

	if currentMonth < 12 {
		currentYear--
	}

	// Generate numbers with active status
	numbers := helpers.GenerateRangeSlice(1, 25)
	days := make([]map[string]interface{}, len(numbers))

	for i, num := range numbers {
		days[i] = map[string]interface{}{
			"Number": num,
			"Active": num <= currentDay,
		}
	}

	data := map[string]interface{}{
		"PageTitle": "Homepage",
		"BodyClass": "home",
		"Days":      days,
		"Month":     currentMonth,
		"Year":      currentYear,
		"Annual":    annual,
	}

	RenderTemplate(w, r, "home", data)
}
