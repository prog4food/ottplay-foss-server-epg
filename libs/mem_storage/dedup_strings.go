package mem_storage

import (
	xxhash32 "github.com/OneOfOne/xxhash"
)

type DedupStrings map[uint32]*string
type DedupByteS map[uint32][]byte


var (
	Str = make(DedupStrings)
	Ids = make(DedupByteS)
	Url = make(DedupStrings)
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
