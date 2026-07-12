package dwg

import (
	"encoding/binary"
	"math"
)

// BitReader reads the DWG bitstream, which is packed MSB-first within each
// byte. Multi-byte raw values (RS, RL, RD) are assembled from consecutive
// bytes in little-endian order.
type BitReader struct {
	data   []byte
	bitPos int // absolute bit offset from start of data
}

func NewBitReader(data []byte) *BitReader {
	return &BitReader{data: data}
}

func (r *BitReader) BitPos() int     { return r.bitPos }
func (r *BitReader) SetBitPos(p int) { r.bitPos = p }
func (r *BitReader) BytePos() int    { return r.bitPos / 8 }
func (r *BitReader) Remaining() int  { return len(r.data)*8 - r.bitPos }
func (r *BitReader) Len() int        { return len(r.data) }

// AlignByte advances to the next byte boundary if not already on one.
func (r *BitReader) AlignByte() {
	if rem := r.bitPos % 8; rem != 0 {
		r.bitPos += 8 - rem
	}
}

// Bit reads a single bit (B).
func (r *BitReader) Bit() int {
	byteIdx := r.bitPos / 8
	if byteIdx >= len(r.data) {
		r.bitPos++
		return 0
	}
	shift := 7 - uint(r.bitPos%8)
	v := (r.data[byteIdx] >> shift) & 1
	r.bitPos++
	return int(v)
}

// Bits reads n bits (n <= 64) MSB-first and returns them right-aligned.
func (r *BitReader) Bits(n int) uint64 {
	var v uint64
	for i := 0; i < n; i++ {
		v = (v << 1) | uint64(r.Bit())
	}
	return v
}

// BB reads a 2-bit code, used as the selector for BS/BL/BD.
func (r *BitReader) BB() int {
	return int(r.Bits(2))
}

// RC reads a raw unsigned char (8 bits).
func (r *BitReader) RC() byte {
	return byte(r.Bits(8))
}

// RS reads a raw short (2 bytes, little-endian byte order).
func (r *BitReader) RS() uint16 {
	b0 := r.RC()
	b1 := r.RC()
	return uint16(b0) | uint16(b1)<<8
}

// RSBE reads a raw short in big-endian byte order, used by the object map.
func (r *BitReader) RSBE() uint16 {
	b0 := r.RC()
	b1 := r.RC()
	return uint16(b0)<<8 | uint16(b1)
}

// RL reads a raw long (4 bytes, little-endian byte order).
func (r *BitReader) RL() uint32 {
	b0 := uint32(r.RC())
	b1 := uint32(r.RC())
	b2 := uint32(r.RC())
	b3 := uint32(r.RC())
	return b0 | b1<<8 | b2<<16 | b3<<24
}

// RD reads a raw IEEE-754 double (8 bytes, little-endian).
func (r *BitReader) RD() float64 {
	lo := uint64(r.RL())
	hi := uint64(r.RL())
	bits := lo | hi<<32
	return math.Float64frombits(bits)
}

// BS reads a bitshort.
func (r *BitReader) BS() uint16 {
	switch r.BB() {
	case 0:
		return r.RS()
	case 1:
		return uint16(r.RC())
	case 2:
		return 0
	default: // 3
		return 256
	}
}

// BSI is BS reinterpreted as a signed 16-bit quantity where useful.
func (r *BitReader) BSI() int16 { return int16(r.BS()) }

// BOT reads a bitcode-object-type (R2010+ object type encoding).
func (r *BitReader) BOT() uint16 {
	switch r.BB() {
	case 0:
		return uint16(r.RC())
	case 1:
		return uint16(r.RC()) + 0x1f0
	default:
		return r.RS()
	}
}

// BL reads a bitlong.
func (r *BitReader) BL() uint32 {
	switch r.BB() {
	case 0:
		return r.RL()
	case 1:
		return uint32(r.RC())
	case 2:
		return 0
	default: // 3, "not used" in spec; treat as 0
		return 0
	}
}

// BLL reads a bitlonglong: a 3-bit length prefix (BB then B) followed by
// that many bytes, least significant first. Used by preview_size on R2010+.
func (r *BitReader) BLL() uint64 {
	n := (r.BB() << 1) | r.Bit()
	var v uint64
	for i := 0; i < n; i++ {
		v |= uint64(r.RC()) << (8 * i)
	}
	return v
}

// BD reads a bitdouble.
func (r *BitReader) BD() float64 {
	switch r.BB() {
	case 0:
		return r.RD()
	case 1:
		return 1.0
	case 2:
		return 0.0
	default:
		return 0.0
	}
}

// BT reads a bitthickness: for R2000+, a flag bit then an optional BD.
func (r *BitReader) BT() float64 {
	if r.Bit() != 0 {
		return 0.0
	}
	return r.BD()
}

// BD2 reads an (x, y) pair of bitdoubles.
func (r *BitReader) BD2() (x, y float64) {
	x = r.BD()
	y = r.BD()
	return
}

// BD3 reads an (x, y, z) triple of bitdoubles.
func (r *BitReader) BD3() (x, y, z float64) {
	x = r.BD()
	y = r.BD()
	z = r.BD()
	return
}

// BE reads a bit-extrusion: defaults to (0,0,1) unless the flag bit is set,
// in which case a raw 3D vector follows as three raw doubles.
func (r *BitReader) BE() (x, y, z float64) {
	if r.Bit() == 0 {
		return 0, 0, 1
	}
	x = r.RD()
	y = r.RD()
	z = r.RD()
	return
}

// DD reads a bitdouble-with-default. The partial forms (BB code 1 and 2)
// splice newly read bytes into specific positions of the default value's
// little-endian byte representation, keeping the rest (mostly sign+exponent)
// from the default - this is how DWG compresses coordinates that are close
// to a previous point.
func (r *BitReader) DD(def float64) float64 {
	switch r.BB() {
	case 0:
		return def
	case 3:
		return r.RD()
	case 1:
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], math.Float64bits(def))
		b[0] = r.RC()
		b[1] = r.RC()
		b[2] = r.RC()
		b[3] = r.RC()
		return math.Float64frombits(binary.LittleEndian.Uint64(b[:]))
	default: // case 2
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], math.Float64bits(def))
		b[4] = r.RC()
		b[5] = r.RC()
		b[0] = r.RC()
		b[1] = r.RC()
		b[2] = r.RC()
		b[3] = r.RC()
		return math.Float64frombits(binary.LittleEndian.Uint64(b[:]))
	}
}

// MC reads a modular char: a variable-length signed varint used for sizes
// and deltas in the object map / handle stream. Little-endian 7-bits-per-byte
// with a continuation flag (0x80); the sign is carried in bit 0x40 of the
// final byte.
func (r *BitReader) MC() int64 {
	var result int64
	shift := uint(0)
	var last byte
	for {
		b := r.RC()
		last = b
		result |= int64(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	if last&0x40 != 0 {
		// clear the sign bit's contribution and negate
		signBitPos := shift
		result &^= int64(0x40) << signBitPos
		result = -result
	}
	return result
}

// MCU reads a modular char as an unsigned value (used for byte counts).
func (r *BitReader) MCU() uint64 {
	var result uint64
	shift := uint(0)
	for {
		b := r.RC()
		result |= uint64(b&0x7f) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return result
}

// MS reads a modular short: 15 data bits + continuation flag per 2-byte
// little-endian short, used by object-map size fields.
func (r *BitReader) MS() uint32 {
	var result uint32
	shift := uint(0)
	for {
		s := r.RS()
		result |= uint32(s&0x7fff) << shift
		if s&0x8000 == 0 {
			break
		}
		shift += 15
	}
	return result
}

// HandleRef is a parsed handle reference (H code).
type HandleRef struct {
	Code  int
	Value uint64 // absolute handle, or delta depending on Code
}

// H reads a handle reference and returns it. It always consumes the right
// number of bytes so callers can skip past fields they don't need.
func (r *BitReader) H() HandleRef {
	first := r.RC()
	code := int(first >> 4)
	counter := int(first & 0x0f)
	var v uint64
	for i := 0; i < counter; i++ {
		v = v<<8 | uint64(r.RC())
	}
	return HandleRef{Code: code, Value: v}
}

// TV reads a text value. useUnicode should be true for DWG versions AC1021
// (2007) and later, which store TV strings as UTF-16LE; earlier versions
// store them as single-byte (ANSI) strings.
func (r *BitReader) TV(useUnicode bool) string {
	n := int(r.BS())
	if n <= 0 {
		return ""
	}
	if !useUnicode {
		b := make([]byte, n)
		for i := range b {
			b[i] = r.RC()
		}
		return string(b)
	}
	u16 := make([]uint16, n)
	for i := range u16 {
		u16[i] = r.RS()
	}
	return decodeUTF16(u16)
}

func decodeUTF16(u []uint16) string {
	runes := make([]rune, 0, len(u))
	for i := 0; i < len(u); i++ {
		r1 := u[i]
		if r1 >= 0xd800 && r1 <= 0xdbff && i+1 < len(u) {
			r2 := u[i+1]
			if r2 >= 0xdc00 && r2 <= 0xdfff {
				runes = append(runes, ((rune(r1)-0xd800)<<10|(rune(r2)-0xdc00))+0x10000)
				i++
				continue
			}
		}
		runes = append(runes, rune(r1))
	}
	return string(runes)
}
