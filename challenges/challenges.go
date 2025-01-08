package challenges

type Challenge struct {
	ID                     int
	Description            string
	Category               string
	MuscleGroups           []string
	Difficulty             int
	CaloriesBurnedEstimate int
}

func GetAllChallenges() []Challenge {
	allChallenges := []Challenge{}
	allChallenges = append(allChallenges, CardioChallenges...)
	allChallenges = append(allChallenges, StrengthChallenges...)
	allChallenges = append(allChallenges, FlexibilityChallenges...)
	allChallenges = append(allChallenges, EnduranceChallenges...)
	return allChallenges
}
