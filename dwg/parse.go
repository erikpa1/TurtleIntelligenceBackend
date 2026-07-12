package dwg

import (
	"fmt"
	"os"
)

// Fixed DWG object type codes for the entities we know how to render.
const (
	objTEXT       = 1
	objARC        = 17
	objCIRCLE     = 18
	objLINE       = 19
	objPOINT      = 27
	objLWPOLYLINE = 77
)

// versionTier captures the bitstream differences between DWG format
// generations that matter for object framing and common entity data. All
// tiers share the same object-map / MS-size-prefixed object layout.
type versionTier struct {
	r2007Plus bool // >= AC1021 (2007): material/shadow flags, separate string stream
	r2010Plus bool // >= AC1024 (2010): BOT object type, UMC handlestream_size, visualstyle flags, BLL preview_size
	r2013Plus bool // >= AC1027 (2013): has_ds_data bit
}

func tierFor(version string) versionTier {
	return versionTier{
		r2007Plus: version >= "AC1021",
		r2010Plus: version >= "AC1024",
		r2013Plus: version >= "AC1027",
	}
}

// ParseFile reads and decodes a DWG file into a Document. Supports AC1018
// (R2004), and the R2007+ family (AC1021/AC1024/AC1027/AC1032, i.e. AutoCAD
// 2007 through 2018+), which share a page/section layout distinct from R2004.
func ParseFile(path string) (*Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse decodes DWG file bytes into a Document containing the flat 2D
// geometry of every LINE, CIRCLE, ARC, POINT, TEXT and LWPOLYLINE entity
// found in the file, regardless of which block owns them.
func Parse(data []byte) (*Document, error) {
	if len(data) < 6 {
		return nil, fmt.Errorf("dwg: file too small")
	}
	version := string(data[0:6])
	switch version {
	// AC1018 (R2004) and R2010+ (AC1024/AC1027/AC1032/...) share the same
	// page/section container format; only object-framing details differ,
	// handled by tierFor. AC1021 (R2007 exactly) uses a different
	// Reed-Solomon-based container and is not wired up (unvalidated).
	case "AC1018":
		return parseR2004(data, version)
	default:
		if version >= "AC1024" {
			return parseR2004(data, version)
		}
		return nil, fmt.Errorf("dwg: unsupported version %q", version)
	}
}

func parseR2004(data []byte, version string) (*Document, error) {
	fh, err := parseFileHeader(data)
	if err != nil {
		return nil, err
	}
	pages, err := parsePageMap(data, fh)
	if err != nil {
		return nil, err
	}
	infos, err := findSectionInfo(data, pages)
	if err != nil {
		return nil, err
	}
	objInfo, ok := infos["AcDb:AcDbObjects"]
	if !ok {
		return nil, fmt.Errorf("dwg: AcDb:AcDbObjects section not found")
	}
	objData, err := readNamedSection(data, pages, objInfo)
	if err != nil {
		return nil, err
	}

	hdlInfo, ok := infos["AcDb:Handles"]
	if !ok {
		return nil, fmt.Errorf("dwg: AcDb:Handles section not found")
	}
	hdlData, err := readNamedSection(data, pages, hdlInfo)
	if err != nil {
		return nil, err
	}

	doc := &Document{Version: version}
	scanObjectMap(objData, hdlData, doc, tierFor(version))
	return doc, nil
}

// scanObjectMap walks the object map (from the decompressed "AcDb:Handles"
// section) to find every object's start offset within the decompressed
// "AcDb:AcDbObjects" stream, then decodes each in turn.
func scanObjectMap(objData, hdlData []byte, doc *Document, tier versionTier) {
	for _, offset := range readObjectMap(hdlData) {
		if offset < 0 || offset+2 > len(objData) {
			continue
		}
		hr := NewBitReader(objData[offset:])
		size := int(hr.MS())
		if tier.r2010Plus {
			hr.MCU() // handlestream_size: unused, we never touch the handle/string streams
		}
		objStart := offset + hr.BytePos()
		if size <= 0 || objStart+size > len(objData) {
			continue
		}
		decodeObjectSafe(objData[objStart:objStart+size], doc, tier)
	}
}

// readObjectMap decodes the "AcDb:Handles" object map: a sequence of pages,
// each holding (handle-delta, offset-delta) pairs that accumulate into the
// absolute byte offset (within the decompressed AcDbObjects stream) of each
// object's own size prefix. Object storage order in AcDbObjects need not
// match handle order, so this map is the only reliable way to find them all.
// The object-map wire format is identical across R2004 and R2007+.
func readObjectMap(data []byte) []int {
	var offsets []int
	r := NewBitReader(data)
	for r.BytePos()+2 <= len(data) {
		// startPos predates the 2-byte size field itself: section_size
		// counts those 2 bytes too, per the reference decoder.
		startPos := r.BytePos()
		sectionSize := int(r.RSBE())
		if sectionSize <= 2 {
			break
		}
		lastOffset := 0 // resets every page: deltas are page-relative, not section-relative
		for r.BytePos()-startPos < sectionSize && r.BytePos() < len(data) {
			oldpos := r.BytePos()
			// handleoff == 0 (or other malformed deltas) is only a warning
			// in the reference decoder, not a stop condition - it still adds
			// the object using the accumulated offset. Only truly making no
			// progress (a real end-of-data condition) should stop the page.
			r.MCU() // handleoff, unused: we don't track handle values
			offset := r.MC()
			if r.BytePos() == oldpos {
				break
			}
			lastOffset += int(offset)
			offsets = append(offsets, lastOffset)
		}
		r.RSBE() // page CRC, unused
		if r.BytePos() >= len(data) {
			break
		}
	}
	return offsets
}

func decodeObjectSafe(objData []byte, doc *Document, tier versionTier) {
	defer func() { recover() }()
	r := NewBitReader(objData)
	decodeObject(r, doc, tier)
}

func decodeObject(r *BitReader, doc *Document, tier versionTier) {
	var objType uint16
	if tier.r2010Plus {
		objType = r.BOT()
	} else {
		objType = r.BS()
	}
	switch objType {
	case objTEXT, objARC, objCIRCLE, objLINE, objPOINT, objLWPOLYLINE:
	default:
		return
	}

	if !tier.r2010Plus {
		r.RL() // bitsize (R2000-R2007): counted separately from the handle stream.
		// R2010+ computes bitsize from handlestream_size instead, so nothing to read here.
	}
	r.H() // the object's own handle
	for { // extended entity data (EED), app-handle-tagged blobs until size 0
		size := r.BS()
		if size == 0 {
			break
		}
		r.H()
		for i := 0; i < int(size); i++ {
			r.RC()
		}
	}
	decodeCommonEntityData(r, tier)

	switch objType {
	case objLINE:
		decodeLine(r, doc)
	case objCIRCLE:
		decodeCircle(r, doc)
	case objARC:
		decodeArc(r, doc)
	case objPOINT:
		decodePoint(r, doc)
	case objTEXT:
		decodeText(r, doc, tier)
	case objLWPOLYLINE:
		decodeLWPolyline(r, doc)
	}
}

// decodeCommonEntityData consumes the AcDbEntity common fields that precede
// every entity's type-specific geometry. Handle fields (layer, ltype,
// material, plotstyle, reactors, xdictionary) live in a separate handle
// stream we never touch, so they cost nothing here.
func decodeCommonEntityData(r *BitReader, tier versionTier) {
	if r.Bit() != 0 { // preview_exists
		var n uint64
		if tier.r2010Plus {
			n = r.BLL()
		} else {
			n = uint64(r.RL())
		}
		// Guard against a garbage/misaligned size turning into a
		// many-billion-iteration loop: it can never legitimately exceed
		// what's left in this object's own byte window.
		if n > uint64(r.Len()) {
			n = 0
		}
		for i := uint64(0); i < n; i++ {
			r.RC()
		}
	}
	r.BB()  // entmode
	r.BL()  // num_reactors
	r.Bit() // is_xdic_missing (SINCE R2004a)
	if tier.r2013Plus {
		r.Bit() // has_ds_data
	}

	flags := r.BS() // color.raw (ENC, SINCE R2004a)
	if flags&0x20 != 0 {
		r.BL() // color.alpha_raw
	}
	if flags&0x40 != 0 {
		// color by handle reference: lives in the handle stream, no bits here
	} else if flags&0x80 != 0 {
		r.BL() // color.rgb
	}
	if flags&0x41 == 0x41 {
		r.TV(false) // color.name
	}
	if flags&0x42 == 0x42 {
		r.TV(false) // color.book_name
	}

	r.BD() // ltype_scale
	r.BB() // ltype_flags
	r.BB() // plotstyle_flags
	if tier.r2007Plus {
		r.BB() // material_flags
		r.RC() // shadow_flags
	}
	if tier.r2010Plus {
		r.Bit() // has_full_visualstyle
		r.Bit() // has_face_visualstyle
		r.Bit() // has_edge_visualstyle
	}
	r.BS() // invisible
	r.RC() // linewt
}

func decodeLine(r *BitReader, doc *Document) {
	zIsZero := r.Bit()
	sx := r.RD()
	ex := r.DD(sx)
	sy := r.RD()
	ey := r.DD(sy)
	var sz, ez float64
	if zIsZero == 0 {
		sz = r.RD()
		ez = r.DD(sz)
	}
	doc.Entities = append(doc.Entities, Entity{
		Kind:   KindLine,
		Points: []Point3{{sx, sy, sz}, {ex, ey, ez}},
	})
}

func decodeCircle(r *BitReader, doc *Document) {
	cx, cy, cz := r.BD3()
	radius := r.BD()
	doc.Entities = append(doc.Entities, Entity{
		Kind:   KindCircle,
		Center: Point3{cx, cy, cz},
		Radius: radius,
	})
}

func decodeArc(r *BitReader, doc *Document) {
	cx, cy, cz := r.BD3()
	radius := r.BD()
	r.BT() // thickness
	r.BE() // extrusion
	startAngle := r.BD()
	endAngle := r.BD()
	doc.Entities = append(doc.Entities, Entity{
		Kind:       KindArc,
		Center:     Point3{cx, cy, cz},
		Radius:     radius,
		StartAngle: startAngle,
		EndAngle:   endAngle,
	})
}

func decodePoint(r *BitReader, doc *Document) {
	x := r.BD()
	y := r.BD()
	z := r.BD()
	doc.Entities = append(doc.Entities, Entity{
		Kind:   KindPoint,
		Points: []Point3{{x, y, z}},
	})
}

// decodeText reads TEXT's geometry fields. On R2007+ the actual string
// content lives in a separate per-object string stream we don't decode, so
// Text is left empty there; the position/height still render correctly.
func decodeText(r *BitReader, doc *Document, tier versionTier) {
	dataflags := r.RC()
	if dataflags&0x01 == 0 {
		r.RD() // elevation
	}
	insX := r.RD()
	insY := r.RD()
	if dataflags&0x02 == 0 {
		r.DD(insX) // alignment_pt.x
		r.DD(insY) // alignment_pt.y
	}
	r.BE() // extrusion
	r.BT() // thickness
	if dataflags&0x04 == 0 {
		r.RD() // oblique_angle
	}
	if dataflags&0x08 == 0 {
		r.RD() // rotation
	}
	height := r.RD()
	if dataflags&0x10 == 0 {
		r.RD() // width_factor
	}
	var text string
	if !tier.r2007Plus {
		text = r.TV(false)
	}
	doc.Entities = append(doc.Entities, Entity{
		Kind:   KindText,
		Points: []Point3{{insX, insY, 0}},
		Text:   text,
		Height: height,
	})
}

func decodeLWPolyline(r *BitReader, doc *Document) {
	flag := r.BS()
	if flag&4 != 0 {
		r.BD() // const_width
	}
	if flag&8 != 0 {
		r.BD() // elevation
	}
	if flag&2 != 0 {
		r.BD() // thickness
	}
	if flag&1 != 0 {
		r.BD3() // extrusion
	}
	numPoints := int(r.BL())
	if numPoints <= 0 || numPoints > 20000 {
		return
	}
	if flag&16 != 0 {
		r.BL() // num_bulges (unused: bulge arcs render as straight segments)
	}
	if flag&32 != 0 {
		r.BL() // num_widths
	}

	points := make([]Point3, numPoints)
	px := r.RD()
	py := r.RD()
	points[0] = Point3{px, py, 0}
	for i := 1; i < numPoints; i++ {
		px = r.DD(px)
		py = r.DD(py)
		points[i] = Point3{px, py, 0}
	}

	doc.Entities = append(doc.Entities, Entity{
		Kind:   KindLWPolyline,
		Points: points,
		Closed: flag&0x200 != 0,
	})
}
