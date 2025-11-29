package levenstein

import (
	"fmt"
)

// LevenshteinDistance calculates the minimum number of single-character edits
// (insertions, deletions, or substitutions) required to change one string into another
func LevenshteinDistance(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)

	// Create a 2D matrix to store distances
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first column (deletions from s1)
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}

	// Initialize first row (insertions to match s2)
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			// Minimum of three operations
			deletion := matrix[i-1][j] + 1
			insertion := matrix[i][j-1] + 1
			substitution := matrix[i-1][j-1] + cost

			matrix[i][j] = min(deletion, min(insertion, substitution))
		}
	}

	return matrix[len1][len2]
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetTestCityNames returns a slice of city names for testing
func GetTestCityNames() []string {
	return []string{
		"New York",
		"Los Angeles",
		"Chicago",
		"Houston",
		"Phoenix",
		"Philadelphia",
		"San Antonio",
		"San Diego",
		"Dallas",
		"San Jose",
		"Austin",
		"Jacksonville",
		"Fort Worth",
		"Columbus",
		"Charlotte",
		"San Francisco",
		"Indianapolis",
		"Seattle",
		"Denver",
		"Boston",
	}
}

// FindClosestMatch finds the closest matching city name using Levenshtein distance
func FindClosestMatch(input string, cities []string) (string, int) {
	if len(cities) == 0 {
		return "", -1
	}

	minDistance := LevenshteinDistance(input, cities[0])
	closestCity := cities[0]

	for i := 1; i < len(cities); i++ {
		distance := LevenshteinDistance(input, cities[i])
		if distance < minDistance {
			minDistance = distance
			closestCity = cities[i]
		}
	}

	return closestCity, minDistance
}

func TestLevenstein() {
	cities := GetTestCityNames()

	fmt.Println("=== Levenshtein Distance Algorithm Demo ===\n")
	fmt.Printf("Available cities: %v\n\n", cities)

	// Test cases with intentional typos
	testInputs := []string{
		"New Yrok",      // Typo in New York
		"Los Angelos",   // Typo in Los Angeles
		"Chikago",       // Typo in Chicago
		"Huston",        // Typo in Houston
		"San Fransisco", // Typo in San Francisco
		"Seatle",        // Typo in Seattle
		"Bostn",         // Typo in Boston
		"Dalla",         // Typo in Dallas
	}

	RunLevensteinAlgorythm(cities, testInputs)

}
