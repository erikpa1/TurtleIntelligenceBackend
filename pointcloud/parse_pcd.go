package pointcloud

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strconv"
	"strings"
)

type pcdField struct {
	name      string
	size      int  // bytes per element
	count     int  // elements in this field (almost always 1)
	typeChar  byte // 'F' float, 'U' unsigned int, 'I' signed int
	rowOffset int  // byte offset of this field's first element within a binary row
	tokenIdx  int  // token index of this field's first element within an ascii line
}

// ParsePCD parses an ASCII or binary PCD (Point Cloud Library) point cloud.
// binary_compressed is not supported.
func ParsePCD(data []byte) ([]Point, bool, error) {
	reader := bufio.NewReader(bytes.NewReader(data))

	var fieldNames, sizeTokens, typeTokens, countTokens []string
	points := 0
	dataMode := ""
	headerBytes := 0

	for {
		line, err := reader.ReadString('\n')
		headerBytes += len(line)
		trimmed := strings.TrimSpace(line)

		if trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			fields := strings.Fields(trimmed)
			key := strings.ToUpper(fields[0])
			switch key {
			case "FIELDS":
				fieldNames = fields[1:]
			case "SIZE":
				sizeTokens = fields[1:]
			case "TYPE":
				typeTokens = fields[1:]
			case "COUNT":
				countTokens = fields[1:]
			case "POINTS":
				if len(fields) >= 2 {
					points, _ = strconv.Atoi(fields[1])
				}
			case "DATA":
				if len(fields) >= 2 {
					dataMode = strings.ToLower(fields[1])
				}
			}
		}

		if strings.HasPrefix(trimmed, "DATA") || err != nil {
			break
		}
	}

	if dataMode == "" || len(fieldNames) == 0 {
		return nil, false, errors.New("invalid pcd header")
	}
	if dataMode == "binary_compressed" {
		return nil, false, errors.New("unsupported pcd data encoding: binary_compressed")
	}

	fields := make([]pcdField, len(fieldNames))
	rowSize := 0
	tokenIdx := 0
	for i, name := range fieldNames {
		size := 4
		if i < len(sizeTokens) {
			size, _ = strconv.Atoi(sizeTokens[i])
		}
		count := 1
		if i < len(countTokens) {
			if c, err := strconv.Atoi(countTokens[i]); err == nil && c > 0 {
				count = c
			}
		}
		typeChar := byte('F')
		if i < len(typeTokens) && len(typeTokens[i]) > 0 {
			typeChar = strings.ToUpper(typeTokens[i])[0]
		}

		fields[i] = pcdField{
			name:      strings.ToLower(name),
			size:      size,
			count:     count,
			typeChar:  typeChar,
			rowOffset: rowSize,
			tokenIdx:  tokenIdx,
		}
		rowSize += size * count
		tokenIdx += count
	}

	xIdx, yIdx, zIdx := -1, -1, -1
	packedIdx := -1
	rIdx, gIdx, bIdx := -1, -1, -1
	for i, f := range fields {
		switch f.name {
		case "x":
			xIdx = i
		case "y":
			yIdx = i
		case "z":
			zIdx = i
		case "rgb", "rgba":
			packedIdx = i
		case "r", "red":
			rIdx = i
		case "g", "green":
			gIdx = i
		case "b", "blue":
			bIdx = i
		}
	}
	if xIdx < 0 || yIdx < 0 || zIdx < 0 {
		return nil, false, errors.New("invalid pcd: missing x/y/z fields")
	}
	hasColor := packedIdx >= 0 || (rIdx >= 0 && gIdx >= 0 && bIdx >= 0)

	result := make([]Point, 0, points)

	// Remaining unread bytes after the header, for binary reading.
	body := data[headerBytes:]

	switch dataMode {
	case "ascii":
		bodyReader := bufio.NewScanner(bytes.NewReader(body))
		bodyReader.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
		read := 0
		for bodyReader.Scan() && read < points {
			line := strings.TrimSpace(bodyReader.Text())
			if line == "" {
				continue
			}
			tokens := strings.Fields(line)
			read++

			x, ex := strconv.ParseFloat(tokens[fields[xIdx].tokenIdx], 64)
			y, ey := strconv.ParseFloat(tokens[fields[yIdx].tokenIdx], 64)
			z, ez := strconv.ParseFloat(tokens[fields[zIdx].tokenIdx], 64)
			if ex != nil || ey != nil || ez != nil {
				continue
			}

			p := Point{X: float32(x), Y: float32(y), Z: float32(z)}
			if hasColor {
				p.R, p.G, p.B = pcdAsciiColor(tokens, fields, packedIdx, rIdx, gIdx, bIdx)
			}
			result = append(result, p)
		}

	case "binary":
		for i := 0; i < points; i++ {
			rowStart := i * rowSize
			if rowStart+rowSize > len(body) {
				break
			}
			row := body[rowStart : rowStart+rowSize]

			x := pcdBinaryScalar(row, fields[xIdx])
			y := pcdBinaryScalar(row, fields[yIdx])
			z := pcdBinaryScalar(row, fields[zIdx])

			p := Point{X: float32(x), Y: float32(y), Z: float32(z)}
			if hasColor {
				p.R, p.G, p.B = pcdBinaryColor(row, fields, packedIdx, rIdx, gIdx, bIdx)
			}
			result = append(result, p)
		}

	default:
		return nil, false, errors.New("unsupported pcd data encoding: " + dataMode)
	}

	if len(result) == 0 {
		return nil, false, ErrNoPoints
	}

	return result, hasColor, nil
}

func pcdAsciiColor(tokens []string, fields []pcdField, packedIdx, rIdx, gIdx, bIdx int) (uint8, uint8, uint8) {
	if packedIdx >= 0 {
		f := fields[packedIdx]
		v, err := strconv.ParseFloat(tokens[f.tokenIdx], 64)
		if err != nil {
			return 0, 0, 0
		}
		var bits uint32
		if f.typeChar == 'F' {
			bits = math.Float32bits(float32(v))
		} else {
			bits = uint32(int64(v))
		}
		return uint8(bits >> 16), uint8(bits >> 8), uint8(bits)
	}

	r, _ := strconv.ParseFloat(tokens[fields[rIdx].tokenIdx], 64)
	g, _ := strconv.ParseFloat(tokens[fields[gIdx].tokenIdx], 64)
	b, _ := strconv.ParseFloat(tokens[fields[bIdx].tokenIdx], 64)
	return normalizeColorTriplet(r, g, b)
}

func pcdBinaryColor(row []byte, fields []pcdField, packedIdx, rIdx, gIdx, bIdx int) (uint8, uint8, uint8) {
	if packedIdx >= 0 {
		f := fields[packedIdx]
		if f.size != 4 || f.rowOffset+4 > len(row) {
			return 0, 0, 0
		}
		bits := binary.LittleEndian.Uint32(row[f.rowOffset:])
		return uint8(bits >> 16), uint8(bits >> 8), uint8(bits)
	}

	r := pcdBinaryScalar(row, fields[rIdx])
	g := pcdBinaryScalar(row, fields[gIdx])
	b := pcdBinaryScalar(row, fields[bIdx])
	return normalizeColorTriplet(r, g, b)
}

func pcdBinaryScalar(row []byte, f pcdField) float64 {
	if f.rowOffset+f.size > len(row) {
		return 0
	}
	raw := row[f.rowOffset:]
	switch {
	case f.typeChar == 'F' && f.size == 4:
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(raw)))
	case f.typeChar == 'F' && f.size == 8:
		return math.Float64frombits(binary.LittleEndian.Uint64(raw))
	case f.size == 1 && f.typeChar == 'I':
		return float64(int8(raw[0]))
	case f.size == 1:
		return float64(raw[0])
	case f.size == 2 && f.typeChar == 'I':
		return float64(int16(binary.LittleEndian.Uint16(raw)))
	case f.size == 2:
		return float64(binary.LittleEndian.Uint16(raw))
	case f.size == 4 && f.typeChar == 'I':
		return float64(int32(binary.LittleEndian.Uint32(raw)))
	case f.size == 4:
		return float64(binary.LittleEndian.Uint32(raw))
	default:
		return 0
	}
}
