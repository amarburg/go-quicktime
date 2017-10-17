package quicktime

import "bytes"
import "encoding/binary"

// int32Decode performs a BigEndian read of an int32 from the data buffer
func int32Decode(data []byte) (ret int32) {
	if len(data) < 4 {
		return 0
	}

	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

// uint32Decode performs a BigEndian read of a uint32 from the data buffer
func uint32Decode(data []byte) (ret uint32) {
	if len(data) < 4 {
		return 0
	}

	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

// fixed32Decdoe reads the "fixed" format used in Quicktime files
// which stores floats as two int16s
func fixed32Decode(data []byte) (ret float32) {
	if len(data) < 4 {
		return 0
	}

	ints := [2]int16{0, 0}
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ints)

	return float32(ints[0]) + float32(ints[1])/65536.0
}
