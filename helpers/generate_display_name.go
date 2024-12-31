package helpers

import (
	"fmt"
	"math/rand"
)

func GenerateDisplayName() string {
	adjectives := []string{
		"Jacked",
		"Swole",
		"Beast",
		"Ripped",
		"Strong",
		"Buff",
		"Shredded",
		"Mighty",
		"Fierce",
		"Massive",
	}

	nouns := []string{
		"Fool",
		"Master",
		"Warrior",
		"Champion",
		"Lifter",
		"Runner",
		"Crusher",
		"Athlete",
		"Dominator",
		"Machine",
	}

	adjective := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	number := rand.Intn(100) + 1

	return fmt.Sprintf("%s%s%d", adjective, noun, number)
}
