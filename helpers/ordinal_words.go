package helpers

import (
	"fmt"
	"strings"
)

func OrdinalWords(num int) string {
	ordinals := map[int]string{
		1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth",
		6: "sixth", 7: "seventh", 8: "eighth", 9: "ninth", 10: "tenth",
		11: "eleventh", 12: "twelfth", 13: "thirteenth", 14: "fourteenth",
		15: "fifteenth", 16: "sixteenth", 17: "seventeenth", 18: "eighteenth",
		19: "nineteenth", 20: "twentieth",
	}

	if word, exists := ordinals[num]; exists {
		return word
	}

	tens := map[int]string{
		2: "twentieth", 3: "thirtieth", 4: "fortieth", 5: "fiftieth",
		6: "sixtieth", 7: "seventieth", 8: "eightieth", 9: "ninetieth",
	}

	baseTens := map[int]string{
		2: "twenty", 3: "thirty", 4: "forty", 5: "fifty",
		6: "sixty", 7: "seventy", 8: "eighty", 9: "ninety",
	}

	if num > 20 && num < 100 {
		tensDigit := num / 10
		unitDigit := num % 10

		if unitDigit == 0 {
			return tens[tensDigit]
		}

		return strings.ToLower(fmt.Sprintf("%s-%s", baseTens[tensDigit], ordinals[unitDigit]))
	}

	return fmt.Sprintf("%dth", num)
}
