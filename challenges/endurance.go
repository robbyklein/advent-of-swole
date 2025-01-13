package challenges

var EnduranceChallenges = []Challenge{
	{
		Description:            "Cycle for 6.2 miles",     // Imperial
		DescriptionMetric:      "Cycle for 10 kilometers", // Metric
		Category:               "Endurance",
		MuscleGroups:           []string{"Legs"},
		Difficulty:             3,
		CaloriesBurnedEstimate: 300,
	},
	{
		Description:            "Swim 546 yards",  // Imperial
		DescriptionMetric:      "Swim 500 meters", // Metric
		Category:               "Endurance",
		MuscleGroups:           []string{"Full Body"},
		Difficulty:             4,
		CaloriesBurnedEstimate: 350,
	},
}
