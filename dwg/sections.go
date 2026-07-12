package dwg

import (
	"encoding/binary"
	"fmt"
)

// fileHeader is the decrypted R2004+ file header (Dwg_R2004_Header in the
// ODA/LibreDWG spec), fields we need to locate the page map and section map.
type fileHeader struct {
	sectionMapID     uint32
	sectionMapAddr   uint64
	sectionArraySize uint32
}

func parseFileHeader(data []byte) (*fileHeader, error) {
	const size = 0x6c
	if len(data) < 0x80+size {
		return nil, fmt.Errorf("dwg: file too small for R2004 header")
	}
	dec := decryptR2004Header(data[0x80 : 0x80+size])
	if string(dec[0:11]) != "AcFssFcAJMB" {
		return nil, fmt.Errorf("dwg: R2004 header decryption failed (bad magic)")
	}
	fh := &fileHeader{
		sectionMapID:     binary.LittleEndian.Uint32(dec[0x50:]),
		sectionMapAddr:   binary.LittleEndian.Uint64(dec[0x54:]),
		sectionArraySize: binary.LittleEndian.Uint32(dec[0x60:]),
	}
	return fh, nil
}

// pageEntry is one physical page's real file address and declared size, as
// found in the (system) Section Page Map.
type pageEntry struct {
	size    uint32
	address uint64
}

// parsePageMap decodes the "Section Page Map" system section, which is
// embedded directly (unencrypted header) right after the file header and
// lists every physical page in the file by number.
func parsePageMap(data []byte, fh *fileHeader) (map[int32]pageEntry, error) {
	base := fh.sectionMapAddr + 0x100
	if base+20 > uint64(len(data)) {
		return nil, fmt.Errorf("dwg: page map address out of range")
	}
	if binary.LittleEndian.Uint32(data[base:]) != 0x41630e3b {
		return nil, fmt.Errorf("dwg: bad page map magic")
	}
	decompSize := binary.LittleEndian.Uint32(data[base+4:])
	compSize := binary.LittleEndian.Uint32(data[base+8:])
	compStart := base + 20
	if compStart+uint64(compSize) > uint64(len(data)) {
		return nil, fmt.Errorf("dwg: page map compressed data out of range")
	}

	dec := make([]byte, decompSize)
	if err := decompressR2004Section(data[compStart:compStart+uint64(compSize)], dec, 0); err != nil {
		return nil, err
	}

	pages := make(map[int32]pageEntry)
	address := uint64(0x100)
	pos := 0
	for pos+8 <= len(dec) {
		number := int32(binary.LittleEndian.Uint32(dec[pos:]))
		size := binary.LittleEndian.Uint32(dec[pos+4:])
		pos += 8
		entryAddr := address
		if number <= int32(fh.sectionArraySize) {
			address += uint64(size)
		}
		if number < 0 {
			pos += 16 // parent/left/right/0x00, only for gap entries
			continue
		}
		pages[number] = pageEntry{size: size, address: entryAddr}
	}

	// The page map's own entry may be listed with a stale cumulative
	// address; the real address is the one we used to find it.
	if e, ok := pages[int32(fh.sectionMapID)]; ok {
		e.address = fh.sectionMapAddr + 0x100
		pages[int32(fh.sectionMapID)] = e
	}
	return pages, nil
}

// sectionInfo describes one named section (e.g. "AcDb:AcDbObjects"): its
// total decompressed size, whether its pages are compressed, and the page
// numbers (in order) that make it up.
type sectionInfo struct {
	name        string
	compressed  uint32
	decompSize  uint64
	pageNumbers []int32
}

// findSectionInfo locates and decodes the "Section Info" system section
// (magic 0x4163003b), which enumerates all named sections in the file. It is
// stored unencrypted, at a page address only found by scanning, since the
// header's section_info_id can be absent/unreliable in the wild.
func findSectionInfo(data []byte, pages map[int32]pageEntry) (map[string]*sectionInfo, error) {
	var base uint64
	found := false
	for _, e := range pages {
		if e.address+4 > uint64(len(data)) {
			continue
		}
		if binary.LittleEndian.Uint32(data[e.address:]) == 0x4163003b {
			base = e.address
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("dwg: section info page not found")
	}

	decompSize := binary.LittleEndian.Uint32(data[base+4:])
	compSize := binary.LittleEndian.Uint32(data[base+8:])
	compStart := base + 20
	if compStart+uint64(compSize) > uint64(len(data)) {
		return nil, fmt.Errorf("dwg: section info compressed data out of range")
	}

	dec := make([]byte, decompSize)
	if err := decompressR2004Section(data[compStart:compStart+uint64(compSize)], dec, 0); err != nil {
		return nil, err
	}

	if len(dec) < 20 {
		return nil, fmt.Errorf("dwg: section info truncated")
	}
	numDesc := binary.LittleEndian.Uint32(dec[0:])
	pos := 20

	result := make(map[string]*sectionInfo)
	for i := uint32(0); i < numDesc; i++ {
		if pos+8+6*4+64 > len(dec) {
			break
		}
		size := binary.LittleEndian.Uint64(dec[pos:])
		pos += 8
		numSections := binary.LittleEndian.Uint32(dec[pos:])
		pos += 4
		pos += 4 // max_decomp_size
		pos += 4 // unknown
		compressed := binary.LittleEndian.Uint32(dec[pos:])
		pos += 4
		pos += 4 // type
		pos += 4 // encrypted
		name := cString(dec[pos : pos+64])
		pos += 64

		si := &sectionInfo{name: name, compressed: compressed, decompSize: size}
		for j := uint32(0); j < numSections; j++ {
			if pos+16 > len(dec) {
				break
			}
			number := int32(binary.LittleEndian.Uint32(dec[pos:]))
			pos += 16 // number(already read) + size(4) + address(8)
			si.pageNumbers = append(si.pageNumbers, number)
		}
		result[name] = si
	}
	return result, nil
}

func cString(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}

// readNamedSection assembles a named section's full decompressed byte
// stream by walking its pages: each page has a 32-byte header XOR-encrypted
// with a mask derived from its own file address, followed by either
// LZ-compressed or raw page data.
func readNamedSection(data []byte, pages map[int32]pageEntry, info *sectionInfo) ([]byte, error) {
	type pageHeader struct {
		address     uint64
		dataSize    uint32
		pageSize    uint32
		startOffset uint64
	}
	var headers []pageHeader
	// Real allocations are sized as num_pages * a per-section-type max
	// (usually 0x7400), not by the section's true logical size - pages can
	// legitimately decompress into more space than that. 0x7400 covers all
	// section types we care about (AcDbObjects, Handles).
	bufSize := info.decompSize
	if want := uint64(len(info.pageNumbers)) * 0x7400; want > bufSize {
		bufSize = want
	}
	for _, num := range info.pageNumbers {
		pe, ok := pages[num]
		if !ok {
			continue
		}
		addr := pe.address
		if addr+32 > uint64(len(data)) {
			continue
		}
		var words [8]uint32
		for k := 0; k < 8; k++ {
			words[k] = binary.LittleEndian.Uint32(data[addr+uint64(k*4):])
		}
		mask := uint32(0x4164536b) ^ uint32(addr)
		for k := 0; k < 8; k++ {
			words[k] ^= mask
		}
		dataSize := words[2]
		pageSize := words[3]
		startOffset := uint64(words[4])
		// Individual pages may be padded/rounded up past the section's
		// nominal decompressed size; size the buffer to fit them all.
		if need := startOffset + uint64(pageSize); need > bufSize {
			bufSize = need
		}
		headers = append(headers, pageHeader{address: addr, dataSize: dataSize, pageSize: pageSize, startOffset: startOffset})
	}

	buf := make([]byte, bufSize)
	for _, h := range headers {
		addr := h.address
		dataSize := h.dataSize
		pageSize := h.pageSize
		startOffset := h.startOffset

		if info.compressed == 2 {
			compStart := addr + 32
			if compStart+uint64(dataSize) > uint64(len(data)) {
				continue
			}
			if startOffset+uint64(pageSize) > uint64(len(buf)) {
				continue
			}
			if err := decompressR2004Section(data[compStart:compStart+uint64(dataSize)], buf, int(startOffset)); err != nil {
				return nil, err
			}
		} else {
			rawStart := addr + 32
			size := uint64(pageSize)
			if startOffset+size > uint64(len(buf)) {
				if uint64(len(buf)) <= startOffset {
					continue
				}
				size = uint64(len(buf)) - startOffset
			}
			if rawStart+size > uint64(len(data)) {
				continue
			}
			copy(buf[startOffset:startOffset+size], data[rawStart:rawStart+size])
		}
	}
	return buf, nil
}
