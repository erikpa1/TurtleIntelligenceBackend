package simMath

import (
	"math"
)

type Position [3]float32

// Distance calculates the Euclidean distance between two positions
func (self Position) Distance(other Position) float32 {
	dx := other[0] - self[0]
	dy := other[1] - self[1]
	dz := other[2] - self[2]
	return float32(math.Sqrt(float64(dx*dx + dy*dy + dz*dz)))
}

// MoveTo moves the position towards destination at the given speed (km/h)
// Returns the distance to the next point (0 if arrived)
func (self *Position) MoveTo(dest Position, speed float32) float32 {
	// Convert speed from km/h to m/s (assuming positions are in meters)
	speedMS := speed / 3.6

	// Calculate distance to destination
	distance := self.Distance(dest)

	// If already at destination, return 0
	if distance < 0.001 {
		return 0
	}

	// Calculate how far we can travel in 1 second at the given speed
	travelDistance := speedMS

	// If we can reach the destination in this step
	if travelDistance >= distance {
		// Move directly to destination
		*self = dest
		return 0
	}

	// Calculate the direction vector (normalized)
	dx := dest[0] - self[0]
	dy := dest[1] - self[1]
	dz := dest[2] - self[2]

	// Normalize the direction vector
	invDistance := 1.0 / distance
	dx *= invDistance
	dy *= invDistance
	dz *= invDistance

	// Move along the direction vector by travelDistance
	self[0] += dx * travelDistance
	self[1] += dy * travelDistance
	self[2] += dz * travelDistance

	// Return remaining distance to destination
	return distance - travelDistance
}
