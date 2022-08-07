package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEloCalc(t *testing.T) {
	dR := calcEloRating(1500, 1, 1500, 0, 40)
	assert.Equal(t, 20, dR)

	dR = calcEloRating(1500, 0, 1500, 1, 40)
	assert.Equal(t, -20, dR)

	dR = calcEloRating(1500, 0.5, 1500, 0.5, 40)
	assert.Equal(t, 0, dR)
}
