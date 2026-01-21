package clay

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

var hasher hash.Hash32 = fnv.New32a()
var hashBuffer []byte = make([]byte, 4+2*4+1)

func hashNumber(offset uint32, seed uint32) ElementId {
	le := binary.LittleEndian
	na := hashBuffer[0:0]
	na = le.AppendUint32(na, offset)
	na = le.AppendUint32(na, seed)

	hasher.Reset()
	hasher.Write(hashBuffer[:2*4])
	hash := hasher.Sum32()

	return ElementId{
		id:       hash + 1, // Reserve the hash result of zero as "null id"
		stringId: STRING_DEFAULT,
	}
}

func hashString(key string) ElementId {
	hasher.Reset()
	hasher.Write([]byte(key))
	hash := hasher.Sum32()

	return ElementId{
		id:       hash + 1, // Reserve the hash result of zero as "null id"
		stringId: key,
	}
}

func hashTextWithConfig(text string, config *TextElementConfig) uint32 {
	hasher.Reset()
	hasher.Write([]byte(text))

	le := binary.LittleEndian
	na := hashBuffer[0:0]
	na = le.AppendUint32(na, uint32(len(text)))
	na = le.AppendUint16(na, config.FontId)
	na = le.AppendUint16(na, config.FontSize)
	na = le.AppendUint16(na, config.LineHeight)
	na = le.AppendUint16(na, config.LetterSpacing)
	na[0] = byte(config.WrapMode)
	hasher.Write(hashBuffer[:4+2*4+1])

	hash := hasher.Sum32()
	return hash + 1 // Reserve the hash result of zero as "null id"
}
