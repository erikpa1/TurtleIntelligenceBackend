package pointcloud

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"
)

// ParseXYZ parses a plain whitespace-separated text point cloud: each
// non-empty line is "x y z" or "x y z r g b". Malformed lines are skipped.
func ParseXYZ(data []byte) ([]Point, bool, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	points := make([]Point, 0, 1024)
	hasColor := false
	colorDecided := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		x, errX := strconv.ParseFloat(fields[0], 64)
		y, errY := strconv.ParseFloat(fields[1], 64)
		z, errZ := strconv.ParseFloat(fields[2], 64)
		if errX != nil || errY != nil || errZ != nil {
			continue
		}

		p := Point{X: float32(x), Y: float32(y), Z: float32(z)}

		if len(fields) >= 6 {
			r, errR := strconv.ParseFloat(fields[3], 64)
			g, errG := strconv.ParseFloat(fields[4], 64)
			b, errB := strconv.ParseFloat(fields[5], 64)
			if errR == nil && errG == nil && errB == nil {
				if !colorDecided {
					hasColor = true
					colorDecided = true
				}
				p.R, p.G, p.B = normalizeColorTriplet(r, g, b)
			}
		} else if !colorDecided {
			colorDecided = true
		}

		points = append(points, p)
	}

	if len(points) == 0 {
		return nil, false, ErrNoPoints
	}

	return points, hasColor, nil
}

// normalizeColorTriplet accepts either 0-255 integer-ish values or 0-1
// normalized floats (common in some xyzrgb exports) and returns 0-255 bytes.
func normalizeColorTriplet(r, g, b float64) (uint8, uint8, uint8) {
	if r <= 1 && g <= 1 && b <= 1 {
		r *= 255
		g *= 255
		b *= 255
	}
	return clampByte(r), clampByte(g), clampByte(b)
}

func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
