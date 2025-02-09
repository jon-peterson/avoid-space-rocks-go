package utils

import "math/rand"

// RndFloat32 returns a random float32 between 0 and max. Panics if <= 0.
func RndFloat32(max float32) float32 {
	if max <= 0 {
		panic("max must be greater than 0")
	}
	return rand.Float32() * max
}

// RndFloat32InRange returns a random float32 between min and max. Panics if min >= max.
func RndFloat32InRange(min, max float32) float32 {
	if min >= max {
		panic("min must be less than max")
	}
	return min + rand.Float32()*(max-min)
}

// RndInt32InRange returns a random int32 between min and max. Panics if min >= max.
func RndInt32InRange(min, max int32) int32 {
	if min >= max {
		panic("min must be less than max")
	}
	return min + rand.Int31n(max-min)
}

// Chance returns true if a random number between 0 and 1 is less than chance.
// Panics if chance is outside the range 0-1.
func Chance(chance float32) bool {
	if chance < 0 || chance > 1 {
		panic("chance must be between 0 and 1")
	}
	return rand.Float32() < chance
}

// Choice returns a random element from the given slice.
func Choice(float32s []float32) float32 {
	return float32s[rand.Intn(len(float32s))]
}
