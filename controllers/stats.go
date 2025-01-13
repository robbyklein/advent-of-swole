package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/robbyklein/swole/config"
	"github.com/robbyklein/swole/db"
	"github.com/robbyklein/swole/sqlc"
)

const (
	pieChartRadius         = 16
	maxDays                = 25
	calorieIncrementsCount = 4
)

type CategoryStat struct {
	Category   string
	Count      int64
	Percentage float64
	Color      string
}

type PieChartData struct {
	Slices []PieSlice
	Legend []CategoryStat
}

type PieSlice struct {
	Path  string
	Color string
}

type caloriesChartData struct {
	Increments []int32
	Bars       []caloriesBar
}

type caloriesBar struct {
	Day        int32
	Calories   int32
	Percentage float64
}

func StatsGET(w http.ResponseWriter, r *http.Request) {
	// Grab the logged in user
	user, ok := r.Context().Value(config.UserContextKey).(sqlc.User)
	if !ok {
		http.Error(w, "Must be logged in", http.StatusInternalServerError)
		return
	}

	// Extract year and month from the URL parameters
	year := chi.URLParam(r, "year")
	month := chi.URLParam(r, "month")

	if year == "" || month == "" {
		http.Error(w, "Year and month are required", http.StatusBadRequest)
		return
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil || yearInt < 2000 || yearInt > 2100 {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	monthInt, err := strconv.Atoi(month)
	if err != nil || monthInt < 1 || monthInt > 12 {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	// Get the challenge month
	challengeMonth, err := db.Queries.GetChallengeMonthByYearMonth(db.CTX, sqlc.GetChallengeMonthByYearMonthParams{
		Year:  int32(yearInt),
		Month: int32(monthInt),
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find challenge month", err), http.StatusInternalServerError)
		return
	}

	// Get calories chart data
	caloriesData, err := getCaloriesChart(user.ID, challengeMonth.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get calories chart data: %v", err), http.StatusInternalServerError)
		return
	}

	// Get muscles pie chart data
	muscleData, err := db.Queries.GetMuscleStats(db.CTX, sqlc.GetMuscleStatsParams{
		ChallengeMonthID: challengeMonth.ID,
		UserID:           user.ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get muscle data: %v", err), http.StatusInternalServerError)
		return
	}
	musclePieChart, err := generateMusclePieChart(muscleData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate muscle pie chart: %v", err), http.StatusInternalServerError)
		return
	}

	// Get categories pie chart data
	categoryData, err := db.Queries.GetCategoryStats(db.CTX, sqlc.GetCategoryStatsParams{
		ChallengeMonthID: challengeMonth.ID,
		UserID:           user.ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get category data: %v", err), http.StatusInternalServerError)
		return
	}
	categoryPieChart, err := generateCategoryPieChart(categoryData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate category pie chart: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch total number of participants
	totalParticipants, err := db.Queries.GetTotalParticipantsForMonth(db.CTX, challengeMonth.ID)
	if err != nil {
		http.Error(w, "Failed to fetch total participants: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch user's rank for the challenge month
	currentPlace, err := db.Queries.GetUserRankForMonth(db.CTX, sqlc.GetUserRankForMonthParams{
		ChallengeMonthID: challengeMonth.ID,
		UserID:           user.ID,
	})
	if err != nil {
		http.Error(w, "Failed to fetch user rank: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch total number of challenges completed by the user
	totalChallengesCompleted, err := db.Queries.GetTotalChallengesCompletedForMonth(db.CTX, sqlc.GetTotalChallengesCompletedForMonthParams{
		ChallengeMonthID: challengeMonth.ID,
		UserID:           user.ID,
	})
	if err != nil {
		http.Error(w, "Failed to fetch total challenges completed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Data passed to the template
	data := map[string]interface{}{
		"PageTitle":                "Stats",
		"Increments":               caloriesData.Increments,
		"Bars":                     caloriesData.Bars,
		"PieChartSlices":           musclePieChart.Slices,
		"Muscles":                  musclePieChart.Legend,
		"CategoryPieChartSlices":   categoryPieChart.Slices,
		"Categories":               categoryPieChart.Legend,
		"TotalParticipants":        totalParticipants,
		"TotalChallengesCompleted": totalChallengesCompleted,
		"CurrentPlace":             currentPlace,
	}

	RenderTemplate(w, r, "stats", data)
}

func getCaloriesChart(userID, ChallengeMonthID int64) (caloriesChartData, error) {
	caloriesRows, err := db.Queries.GetCaloriesStats(db.CTX, sqlc.GetCaloriesStatsParams{
		ChallengeMonthID: ChallengeMonthID,
		UserID:           userID,
	})
	if err != nil {
		return caloriesChartData{}, err
	}

	caloriesMap := make(map[int32]int32)
	var maxDay int32

	for _, row := range caloriesRows {
		totalCalories, ok := row.TotalCalories.(int32)
		if !ok {
			return caloriesChartData{}, errors.New("invalid data format for TotalCalories")
		}
		caloriesMap[row.DayNumber] = totalCalories
		if totalCalories > maxDay {
			maxDay = totalCalories
		}
	}

	// Calculate increments for the left side of the chart
	increment := int32(math.Ceil(float64(maxDay) / float64(calorieIncrementsCount)))
	increments := []int32{0}
	for i := 1; i < calorieIncrementsCount; i++ {
		increments = append(increments, increment*int32(i))
	}
	increments = append(increments, maxDay)

	// Generate bars, including placeholders for missing days (up to maxDays)
	bars := make([]caloriesBar, 0, maxDays)
	for day := int32(1); day <= maxDays; day++ {
		totalCalories := caloriesMap[day] // Default to 0 if not in map
		percentage := 0.0
		if maxDay > 0 {
			percentage = float64(totalCalories) / float64(maxDay) * 100
		}
		bars = append(bars, caloriesBar{
			Day:        day,
			Calories:   totalCalories,
			Percentage: percentage,
		})
	}

	return caloriesChartData{
		Increments: increments,
		Bars:       bars,
	}, nil
}

func generateMusclePieChart(rows []sqlc.GetMuscleStatsRow) (PieChartData, error) {
	// Special check: if there's exactly 1 row and it's 100% (or close), draw a full circle
	if len(rows) == 1 {
		r := rows[0]
		if r.Percentage > 99.9999 {
			color := generateColor(0)
			legend := []CategoryStat{{
				Category:   fmt.Sprintf("%v", r.MuscleGroup),
				Count:      r.Count,
				Percentage: r.Percentage,
				Color:      color,
			}}

			// This path describes a full circle of radius 16 centered at (16,16).
			fullCirclePath := "M 16 16 m -16 0 a 16 16 0 1 0 32 0 a 16 16 0 1 0 -32 0"

			return PieChartData{
				Slices: []PieSlice{{
					Path:  fullCirclePath,
					Color: color,
				}},
				Legend: legend,
			}, nil
		}
	}

	legend := make([]CategoryStat, 0, len(rows))
	pieSlices := make([]PieSlice, 0, len(rows))
	currentAngle := 0.0

	for i, row := range rows {
		group := fmt.Sprintf("%v", row.MuscleGroup)
		color := generateColor(i)

		legend = append(legend, CategoryStat{
			Category:   group,
			Count:      row.Count,
			Percentage: row.Percentage,
			Color:      color,
		})

		if row.Percentage == 0 {
			continue
		}

		sliceAngle := (row.Percentage / 100.0) * 360.0
		largeArcFlag := 0
		if sliceAngle > 180 {
			largeArcFlag = 1
		}

		startX := pieChartRadius + math.Cos(currentAngle*(math.Pi/180))*pieChartRadius
		startY := pieChartRadius + math.Sin(currentAngle*(math.Pi/180))*pieChartRadius

		endX := pieChartRadius + math.Cos((currentAngle+sliceAngle)*(math.Pi/180))*pieChartRadius
		endY := pieChartRadius + math.Sin((currentAngle+sliceAngle)*(math.Pi/180))*pieChartRadius

		path := fmt.Sprintf(
			"M %d %d L %.2f %.2f A %d %d 0 %d 1 %.2f %.2f Z",
			pieChartRadius, pieChartRadius,
			startX, startY,
			pieChartRadius, pieChartRadius,
			largeArcFlag,
			endX, endY,
		)

		pieSlices = append(pieSlices, PieSlice{
			Path:  path,
			Color: color,
		})

		currentAngle += sliceAngle
	}

	return PieChartData{
		Slices: pieSlices,
		Legend: legend,
	}, nil
}

func generateCategoryPieChart(rows []sqlc.GetCategoryStatsRow) (PieChartData, error) {
	if len(rows) == 1 {
		r := rows[0]
		if r.Percentage > 99.9999 {
			color := generateColor(0)
			legend := []CategoryStat{{
				Category:   r.Category,
				Count:      r.Count,
				Percentage: r.Percentage,
				Color:      color,
			}}

			fullCirclePath := "M 16 16 m -16 0 a 16 16 0 1 0 32 0 a 16 16 0 1 0 -32 0"

			return PieChartData{
				Slices: []PieSlice{{
					Path:  fullCirclePath,
					Color: color,
				}},
				Legend: legend,
			}, nil
		}
	}

	legend := make([]CategoryStat, 0, len(rows))
	pieSlices := make([]PieSlice, 0, len(rows))
	currentAngle := 0.0

	for i, row := range rows {
		color := generateColor(i)

		legend = append(legend, CategoryStat{
			Category:   row.Category,
			Count:      row.Count,
			Percentage: row.Percentage,
			Color:      color,
		})

		if row.Percentage == 0 {
			continue
		}

		sliceAngle := (row.Percentage / 100.0) * 360.0
		largeArcFlag := 0
		if sliceAngle > 180 {
			largeArcFlag = 1
		}

		startX := pieChartRadius + math.Cos(currentAngle*(math.Pi/180))*pieChartRadius
		startY := pieChartRadius + math.Sin(currentAngle*(math.Pi/180))*pieChartRadius

		endX := pieChartRadius + math.Cos((currentAngle+sliceAngle)*(math.Pi/180))*pieChartRadius
		endY := pieChartRadius + math.Sin((currentAngle+sliceAngle)*(math.Pi/180))*pieChartRadius

		path := fmt.Sprintf(
			"M %d %d L %.2f %.2f A %d %d 0 %d 1 %.2f %.2f Z",
			pieChartRadius, pieChartRadius,
			startX, startY,
			pieChartRadius, pieChartRadius,
			largeArcFlag,
			endX, endY,
		)

		pieSlices = append(pieSlices, PieSlice{
			Path:  path,
			Color: color,
		})

		currentAngle += sliceAngle
	}

	return PieChartData{
		Slices: pieSlices,
		Legend: legend,
	}, nil
}

func generateColor(index int) string {
	colors := []string{
		"#2a73ef",
		"#ef2a2a",
		"#2aef93",
		"#efe22a",
		"#ef902a",
		"#682aef",
		"#89ef2a",
		"#ef2aa7",
		"#c42aef",
	}
	return colors[index%len(colors)]
}
