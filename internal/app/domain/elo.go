package domain

import "math"

func calcEloRating(Ra, Rb, Sa, Sb, K int) (int, int) {
	var Ea, Eb int
	Ea = int(1 / (1 + math.Pow(float64(10), float64((Rb-Ra)/400))))
	Eb = int(1 / (1 + math.Pow(float64(10), float64((Ra-Rb)/400))))

	return K * (Sa - Ea), K * (Sb - Eb)
}
