package dwg

import "math"

type Point3 struct {
	X, Y, Z float64
}

type EntityKind int

const (
	KindLine EntityKind = iota
	KindCircle
	KindArc
	KindPoint
	KindText
	KindLWPolyline
)

func (k EntityKind) String() string {
	switch k {
	case KindLine:
		return "LINE"
	case KindCircle:
		return "CIRCLE"
	case KindArc:
		return "ARC"
	case KindPoint:
		return "POINT"
	case KindText:
		return "TEXT"
	case KindLWPolyline:
		return "LWPOLYLINE"
	default:
		return "UNKNOWN"
	}
}

// Entity is a decoded, renderable DWG entity. Not every DWG field is kept -
// only what's needed to rasterize basic 2D geometry.
type Entity struct {
	Kind       EntityKind
	Handle     uint64
	Points     []Point3 // LINE: 2 pts; POINT/TEXT: 1 pt; LWPOLYLINE: N vertices
	Center     Point3   // CIRCLE/ARC
	Radius     float64  // CIRCLE/ARC
	StartAngle float64  // ARC, radians
	EndAngle   float64  // ARC, radians
	Closed     bool     // LWPOLYLINE
	Text       string   // TEXT
	Height     float64  // TEXT
}

// Document holds everything decoded from a DWG file that's needed to render
// a flat 2D raster preview.
type Document struct {
	Version  string
	Entities []Entity
}

// BoundingBox returns the min/max corners across all entities that
// contribute visible geometry. ok is false if there's nothing to bound.
func (d *Document) BoundingBox() (min, max Point3, ok bool) {
	first := true
	consider := func(p Point3) {
		if first {
			min, max = p, p
			first = false
			return
		}
		if p.X < min.X {
			min.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}
	for _, e := range d.Entities {
		switch e.Kind {
		case KindCircle:
			consider(Point3{e.Center.X - e.Radius, e.Center.Y - e.Radius, 0})
			consider(Point3{e.Center.X + e.Radius, e.Center.Y + e.Radius, 0})
		case KindArc:
			// conservative: bound by full circle, arcs are a subset
			consider(Point3{e.Center.X - e.Radius, e.Center.Y - e.Radius, 0})
			consider(Point3{e.Center.X + e.Radius, e.Center.Y + e.Radius, 0})
		default:
			for _, p := range e.Points {
				consider(p)
			}
		}
	}
	return min, max, !first
}

// Polyline tessellates an entity into one or more open/closed polylines of
// 2D points, suitable for line-segment rasterization. Circles and arcs are
// sampled into straight segments.
func (e Entity) Polylines() [][]Point3 {
	const arcSteps = 64
	switch e.Kind {
	case KindLine:
		return [][]Point3{e.Points}
	case KindPoint:
		if len(e.Points) == 1 {
			p := e.Points[0]
			d := 0.0001
			return [][]Point3{{{p.X - d, p.Y, 0}, {p.X + d, p.Y, 0}}}
		}
	case KindLWPolyline:
		pts := e.Points
		if e.Closed && len(pts) > 1 {
			pts = append(append([]Point3{}, pts...), pts[0])
		}
		return [][]Point3{pts}
	case KindCircle:
		pts := make([]Point3, 0, arcSteps+1)
		for i := 0; i <= arcSteps; i++ {
			a := 2 * math.Pi * float64(i) / float64(arcSteps)
			pts = append(pts, Point3{
				X: e.Center.X + e.Radius*math.Cos(a),
				Y: e.Center.Y + e.Radius*math.Sin(a),
			})
		}
		return [][]Point3{pts}
	case KindArc:
		start, end := e.StartAngle, e.EndAngle
		for end < start {
			end += 2 * math.Pi
		}
		pts := make([]Point3, 0, arcSteps+1)
		for i := 0; i <= arcSteps; i++ {
			a := start + (end-start)*float64(i)/float64(arcSteps)
			pts = append(pts, Point3{
				X: e.Center.X + e.Radius*math.Cos(a),
				Y: e.Center.Y + e.Radius*math.Sin(a),
			})
		}
		return [][]Point3{pts}
	}
	return nil
}
