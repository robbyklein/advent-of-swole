package challenges

var FlexibilityChallenges = []Challenge{
	{
		Description:            "Hold a plank for 2 minutes", // Same for both units (time-based)
		DescriptionMetric:      "Hold a plank for 2 minutes",
		Category:               "Flexibility",
		MuscleGroups:           []string{"Core", "Shoulders"},
		Difficulty:             2,
		CaloriesBurnedEstimate: 20,
	},
	{
		Description:            "Do 15 minutes of yoga", // Same for both units (time-based)
		DescriptionMetric:      "Do 15 minutes of yoga",
		Category:               "Flexibility",
		MuscleGroups:           []string{"Full Body"},
		Difficulty:             2,
		CaloriesBurnedEstimate: 50,
	},
}
