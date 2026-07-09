package pointcloud

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"strings"
)

type plyProperty struct {
	name string
	size int // bytes for scalar types; list properties are not supported for vertex data
}

func plyTypeSize(t string) int {
	switch t {
	case "char", "uchar", "int8", "uint8":
		return 1
	case "short", "ushort", "int16", "uint16":
		return 2
	case "int", "uint", "int32", "uint32", "float", "float32":
		return 4
	case "double", "float64":
		return 8
	default:
		return 0
	}
}

// ParsePLY parses an ASCII or binary (little/big endian) PLY point cloud.
// Only the "vertex" element is read; any other elements (e.g. faces) are
// ignored.
func ParsePLY(data []byte) ([]Point, bool, error) {
	headerEnd := bytes.Index(data, []byte("end_header"))
	if headerEnd < 0 {
		return nil, false, errors.New("invalid ply: missing end_header")
	}

	header := string(data[:headerEnd])
	// Body starts right after the "end_header" line's newline.
	bodyStart := headerEnd + len("end_header")
	for bodyStart < len(data) && (data[bodyStart] == '\r' || data[bodyStart] == '\n') {
		bodyStart++
	}
	body := data[bodyStart:]

	format := ""
	vertexCount := 0
	var vertexProps []plyProperty
	inVertexElement := false

	for _, rawLine := range strings.Split(header, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || line == "ply" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "format":
			if len(fields) >= 2 {
				format = fields[1]
			}
		case "comment", "obj_info":
			continue
		case "element":
			if len(fields) >= 3 {
				inVertexElement = fields[1] == "vertex"
				if inVertexElement {
					vertexCount, _ = strconv.Atoi(fields[2])
				}
			}
		case "property":
			if !inVertexElement {
				continue
			}
			if fields[1] == "list" {
				// Not supported/expected on vertex elements for point clouds; skip.
				continue
			}
			size := plyTypeSize(fields[1])
			if size == 0 {
				continue
			}
			vertexProps = append(vertexProps, plyProperty{name: strings.ToLower(fields[2]), size: size})
		}
	}

	if vertexCount == 0 || len(vertexProps) == 0 {
		return nil, false, errors.New("invalid ply: no vertex element")
	}

	xIdx, yIdx, zIdx := -1, -1, -1
	rIdx, gIdx, bIdx := -1, -1, -1
	offsets := make([]int, len(vertexProps))
	offset := 0
	for i, p := range vertexProps {
		offsets[i] = offset
		offset += p.size
		switch p.name {
		case "x":
			xIdx = i
		case "y":
			yIdx = i
		case "z":
			zIdx = i
		case "red", "r", "diffuse_red":
			rIdx = i
		case "green", "g", "diffuse_green":
			gIdx = i
		case "blue", "b", "diffuse_blue":
			bIdx = i
		}
	}
	rowSize := offset

	if xIdx < 0 || yIdx < 0 || zIdx < 0 {
		return nil, false, errors.New("invalid ply: missing x/y/z properties")
	}
	hasColor := rIdx >= 0 && gIdx >= 0 && bIdx >= 0

	points := make([]Point, 0, vertexCount)

	// Re-derive property type names for scalar decode (need actual PLY type
	// strings, not just sizes) - re-parse header just for vertex property types.
	propTypes := plyPropertyTypes(header)

	switch format {
	case "ascii":
		lines := strings.Split(string(body), "\n")
		read := 0
		for _, rawLine := range lines {
			if read >= vertexCount {
				break
			}
			line := strings.TrimSpace(rawLine)
			if line == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < len(vertexProps) {
				continue
			}

			x, ex := strconv.ParseFloat(fields[xIdx], 64)
			y, ey := strconv.ParseFloat(fields[yIdx], 64)
			z, ez := strconv.ParseFloat(fields[zIdx], 64)
			if ex != nil || ey != nil || ez != nil {
				read++
				continue
			}
			p := Point{X: float32(x), Y: float32(y), Z: float32(z)}
			if hasColor {
				r, _ := strconv.ParseFloat(fields[rIdx], 64)
				g, _ := strconv.ParseFloat(fields[gIdx], 64)
				b, _ := strconv.ParseFloat(fields[bIdx], 64)
				p.R, p.G, p.B = normalizeColorTriplet(r, g, b)
			}
			points = append(points, p)
			read++
		}

	case "binary_little_endian", "binary_big_endian":
		var order binary.ByteOrder = binary.LittleEndian
		if format == "binary_big_endian" {
			order = binary.BigEndian
		}

		for i := 0; i < vertexCount; i++ {
			rowStart := i * rowSize
			if rowStart+rowSize > len(body) {
				break
			}
			row := body[rowStart : rowStart+rowSize]

			x := readPlyScalar(row, offsets[xIdx], propTypes[xIdx], order)
			y := readPlyScalar(row, offsets[yIdx], propTypes[yIdx], order)
			z := readPlyScalar(row, offsets[zIdx], propTypes[zIdx], order)

			p := Point{X: float32(x), Y: float32(y), Z: float32(z)}
			if hasColor {
				r := readPlyScalar(row, offsets[rIdx], propTypes[rIdx], order)
				g := readPlyScalar(row, offsets[gIdx], propTypes[gIdx], order)
				b := readPlyScalar(row, offsets[bIdx], propTypes[bIdx], order)
				p.R, p.G, p.B = normalizeColorTriplet(r, g, b)
			}
			points = append(points, p)
		}

	default:
		return nil, false, errors.New("unsupported ply format: " + format)
	}

	if len(points) == 0 {
		return nil, false, ErrNoPoints
	}

	return points, hasColor, nil
}

// plyPropertyTypes re-walks the header text to recover the PLY type string
// (e.g. "float", "uchar") for each vertex property, in declaration order.
func plyPropertyTypes(header string) []string {
	var types []string
	inVertexElement := false
	for _, rawLine := range strings.Split(header, "\n") {
		fields := strings.Fields(strings.TrimSpace(rawLine))
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "element":
			inVertexElement = len(fields) >= 2 && fields[1] == "vertex"
		case "property":
			if inVertexElement && len(fields) >= 3 && fields[1] != "list" {
				types = append(types, fields[1])
			}
		}
	}
	return types
}

func readPlyScalar(row []byte, offset int, plyType string, order binary.ByteOrder) float64 {
	switch plyType {
	case "char", "int8":
		return float64(int8(row[offset]))
	case "uchar", "uint8":
		return float64(row[offset])
	case "short", "int16":
		return float64(int16(order.Uint16(row[offset:])))
	case "ushort", "uint16":
		return float64(order.Uint16(row[offset:]))
	case "int", "int32":
		return float64(int32(order.Uint32(row[offset:])))
	case "uint", "uint32":
		return float64(order.Uint32(row[offset:]))
	case "float", "float32":
		return float64(math.Float32frombits(order.Uint32(row[offset:])))
	case "double", "float64":
		return math.Float64frombits(order.Uint64(row[offset:]))
	default:
		return 0
	}
}
