package main

import (
	"math"
	"math/rand"
)

func sampleLaplace(b float64) float64 {
	// Generate a uniform random number in [0,1)
	u := rand.Float64()
	// Use the inverse CDF method based on u:
	if u < 0.5 {
		return b * math.Log(2*u)
	}
	return -b * math.Log(2*(1-u))
}

func addNoise(trueValue float64, epsilon, sensitivity float64) float64 {
	// Calculate the scale parameter for the Laplace distribution
	b := sensitivity / epsilon

	// Sample from the Laplace distribution
	noise := sampleLaplace(b)

	// Add noise to the true value
	return trueValue + noise
}
