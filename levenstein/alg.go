package levenstein

import "fmt"

func RunLevensteinAlgorythm(sourceStrings []string, searched []string) {

	fmt.Println("=== Finding Closest Matches for Misspelled Cities ===\n")
	for _, input := range searched {
		closestCity, distance := FindClosestMatch(input, sourceStrings)
		fmt.Printf("Input: %-15s -> Closest Match: %-15s (Distance: %d)\n",
			input, closestCity, distance)
	}

	// Calculate distances between all pairs of selected cities
	fmt.Println("\n=== Distance Matrix for Selected Cities ===\n")
	selectedCities := []string{"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"}

	fmt.Printf("%-15s", "")
	for _, city := range selectedCities {
		fmt.Printf("%-12s", city)
	}
	fmt.Println()

	for _, city1 := range selectedCities {
		fmt.Printf("%-15s", city1)
		for _, city2 := range selectedCities {
			distance := LevenshteinDistance(city1, city2)
			fmt.Printf("%-12d", distance)
		}
		fmt.Println()
	}

	// Example of exact distance calculations
	fmt.Println("\n=== Detailed Distance Examples ===\n")
	examples := [][]string{
		{"cat", "hat"},
		{"Saturday", "Sunday"},
		{"kitten", "sitting"},
		{"book", "back"},
		{"algorithm", "altruistic"},
	}

	for _, pair := range examples {
		distance := LevenshteinDistance(pair[0], pair[1])
		fmt.Printf("'%s' -> '%s': Distance = %d\n", pair[0], pair[1], distance)
	}

}
