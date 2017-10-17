package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "fmt"

// A CO64Atom is a table of 64bit chunk offsets
// For simplicity, I internally convert STCO atoms to CO64, rather than
// duplicating functionality
type CO64Atom struct {
	Atom         *Atom
	ChunkOffsets []int64
}

// ParseSTCO converts a generic "stco" or "co64" atom to a CO64Atom struct
func ParseSTCO(atom *Atom) (CO64Atom, error) {

	if !atom.HasData() {
		return CO64Atom{}, errors.New("STCO Atom doesn't have data")
	}

	// TODO:  Fix the DRY
	if atom.Type == "stco" {

		numEntries := uint32Decode(atom.Data[4:8])

		int32Offsets := make([]int32, numEntries)
		buf := bytes.NewBuffer(atom.Data[8:])
		binary.Read(buf, binary.BigEndian, &int32Offsets)

		stco := CO64Atom{Atom: atom,
			ChunkOffsets: make([]int64, numEntries)}

		// Copy the int32s to int64s
		// TODO:  I'm sure this isn't idiomatic
		for i := 0; i < int(numEntries); i++ {
			stco.ChunkOffsets[i] = int64(int32Offsets[i])
		}

		return stco, nil

	} else if atom.Type == "co64" {

		numEntries := uint32Decode(atom.Data[4:8])

		stco := CO64Atom{Atom: atom,
			ChunkOffsets: make([]int64, numEntries, numEntries)}

		buf := bytes.NewBuffer(atom.Data[8:])
		for i := 0; i < int(numEntries); i++ {
			binary.Read(buf, binary.BigEndian, &stco.ChunkOffsets[i])
		}

		if len(stco.ChunkOffsets) != int(numEntries) {
			panic(fmt.Sprintf("Unexpected state: len(stco.ChunkOffsets) != numEntries, %d != %d", len(stco.ChunkOffsets), numEntries))
		}

		return stco, nil

	} else {
		return CO64Atom{}, errors.New("Not an STCO atom")
	}
}

// ChunkOffset gives the offset for a given chunk (remember that chunks are
// enumerated base-1)
func (stco CO64Atom) ChunkOffset(chunk int) int64 {
	if chunk < 1 || chunk > len(stco.ChunkOffsets) {
		panic(fmt.Sprintf("Requested chunk %d in file with %d chunks", chunk, len(stco.ChunkOffsets)))
	}

	return stco.ChunkOffsets[chunk-1]
}
