package compressservice

import (
	"github.com/jackma8ge8/pine/serializer"
)

type serverCompress struct {
	kindToCode map[string]byte
	codeToKind map[byte]string
}

// AddRecord add serverKind and serverCode recore
func (sc *serverCompress) AddRecord(serverKind string) {
	if _, exist := sc.kindToCode[serverKind]; exist {
		return
	}

	code := byte(len(sc.kindToCode) + 1)
	sc.kindToCode[serverKind] = code
	sc.codeToKind[code] = serverKind
}

// GetKindByCode get serverKind by serverCode
func (sc *serverCompress) GetKindByCode(code byte) string {
	if value, exist := sc.codeToKind[code]; exist {
		return value
	}
	return ""
}

// GetCodeByKind get serverCode by serverKind
func (sc *serverCompress) GetCodeByKind(serverKind string) byte {
	if value, exist := sc.kindToCode[serverKind]; exist {
		return value
	}
	return 0
}

// ToBytes get json bytes
func (sc *serverCompress) ToBytes() []byte {

	return serializer.ToBytes(map[string]interface{}{
		"kindToCode": sc.kindToCode,
		"codeToKind": sc.codeToKind,
	})
}

// Server map
var Server = serverCompress{
	kindToCode: make(map[string]byte),
	codeToKind: make(map[byte]string),
}
