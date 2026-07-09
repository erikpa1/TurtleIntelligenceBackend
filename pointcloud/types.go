package pointcloud

import (
	"errors"
	"math"
)

// Point is a single parsed point cloud sample. Color is only meaningful when
// the source file carried color data (see Parse's hasColor return value).
type Point struct {
	X, Y, Z float32
	R, G, B uint8
}

type Bounds struct {
	Min [3]float64
	Max [3]float64
}

func BoundsOf(points []Point) Bounds {
	b := Bounds{
		Min: [3]float64{math.Inf(1), math.Inf(1), math.Inf(1)},
		Max: [3]float64{math.Inf(-1), math.Inf(-1), math.Inf(-1)},
	}
	for _, p := range points {
		v := [3]float64{float64(p.X), float64(p.Y), float64(p.Z)}
		for i := 0; i < 3; i++ {
			if v[i] < b.Min[i] {
				b.Min[i] = v[i]
			}
			if v[i] > b.Max[i] {
				b.Max[i] = v[i]
			}
		}
	}
	return b
}

func (b Bounds) Center() [3]float64 {
	return [3]float64{
		(b.Min[0] + b.Max[0]) / 2,
		(b.Min[1] + b.Max[1]) / 2,
		(b.Min[2] + b.Max[2]) / 2,
	}
}

// Octant returns the bounding box of child octant i (0-7) of this bounds,
// split at center. Bit 0 of i selects the X half, bit 1 the Y half, bit 2
// the Z half (0 = lower half, 1 = upper half).
func (b Bounds) Octant(i int, center [3]float64) Bounds {
	child := Bounds{Min: b.Min, Max: b.Max}
	for axis := 0; axis < 3; axis++ {
		if i&(1<<axis) != 0 {
			child.Min[axis] = center[axis]
		} else {
			child.Max[axis] = center[axis]
		}
	}
	return child
}

var ErrUnsupportedFormat = errors.New("unsupported point cloud format")
var ErrNoPoints = errors.New("no points parsed from file")

// Parse dispatches to the parser matching the given file extension (without
// the leading dot, case-insensitive).
func Parse(ext string, data []byte) ([]Point, bool, error) {
	switch normalizeExt(ext) {
	case "ply":
		return ParsePLY(data)
	case "pcd":
		return ParsePCD(data)
	case "xyz", "txt", "xyzrgb", "xyzn":
		return ParseXYZ(data)
	default:
		return nil, false, ErrUnsupportedFormat
	}
}

func normalizeExt(ext string) string {
	out := make([]byte, 0, len(ext))
	for _, c := range ext {
		if c >= 'A' && c <= 'Z' {
			c = c - 'A' + 'a'
		}
		out = append(out, byte(c))
	}
	return string(out)
}
