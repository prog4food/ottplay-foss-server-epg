package mem_storage

import (
	xxhash32 "github.com/OneOfOne/xxhash"

	"ottplay-foss-server-epg/pkg/helpers"
)

type DedupStrings map[uint32]*string
type DedupByteS map[uint32][]byte


var (
	// Str = make(DedupStrings)
	Ids = make(DedupByteS)
	Url = make(DedupByteS)
)

func DedupStr(t DedupStrings, s *string) *string {
	_h := xxhash32.ChecksumString32(*s)
	if val, ok := t[_h]; ok {
		return val
	}
 	t[_h] = s
 	return s
}


func DedupByte(t DedupByteS, s []byte) []byte {
	_h := xxhash32.Checksum32(s)
	if val, ok := t[_h]; ok {
		return val
	}
 	t[_h] = s
 	return s
}

// MemUnsafe
func DedupByteByS(t DedupByteS, s *string) []byte {
	_b := helpers.S2bP(s)
	_h := xxhash32.Checksum32(_b)
	if val, ok := t[_h]; ok {
		return val
	}
 	t[_h] = _b
 	return _b
}
