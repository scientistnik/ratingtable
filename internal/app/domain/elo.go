package domain

import "math"

func calcEloRating(Ra int, Sa float64, Rb int, Sb float64, K int) int {
	return int(float64(K) * (Sa - (1 / (1 + math.Pow(float64(10), float64((Rb-Ra)/400))))))
}
