package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "fmt"

// I just cast everything to a CO64 atom for in-memory representation
type CO64Atom struct {
	Atom         *Atom
	ChunkOffsets []int64
}

// Also handle CO64
func ParseSTCO(atom *Atom) (CO64Atom, error) {

	if !atom.HasData() {
		return CO64Atom{}, errors.New("STCO Atom doesn't have data")
	}

	// TODO:  Fix the DRY
	if atom.Type == "stco" {

		num_entries := Uint32Decode(atom.Data[4:8])

		int32Offsets := make([]int32, num_entries)
		buf := bytes.NewBuffer(atom.Data[8:])
		binary.Read(buf, binary.BigEndian, &int32Offsets)

		stco := CO64Atom{Atom: atom,
			ChunkOffsets: make([]int64, num_entries)}

		// Copy the int32s to int64s
		// TODO:  I'm sure this isn't idiomatic
		for i := 0; i < int(num_entries); i++ {
			stco.ChunkOffsets[i] = int64(int32Offsets[i])
		}

		return stco, nil

	} else if atom.Type == "co64" {

		numEntries := Uint32Decode(atom.Data[4:8])
//fmt.Printf("co64 has %d entries\n", numEntries )

		stco := CO64Atom{Atom: atom,
										ChunkOffsets: make([]int64, numEntries, numEntries) }

		buf := bytes.NewBuffer(atom.Data[8:])
		for i := 0; i < int(numEntries); i++ {
			binary.Read(buf, binary.BigEndian, &stco.ChunkOffsets[i])
		}

		if len(stco.ChunkOffsets) != int(numEntries) {
			panic(fmt.Sprintf("Unexpected state: len(stco.ChunkOffsets) != numEntries, %d != %d",len(stco.ChunkOffsets),numEntries))
		}

		return stco, nil

	} else {
		return CO64Atom{}, errors.New("Not an STCO atom")
	}
}

// Remember that chunks are 1-base
func (stco CO64Atom) ChunkOffset(chunk int) int64 {
	if chunk < 1 || chunk > len(stco.ChunkOffsets) {
		panic(fmt.Sprintf("Requested chunk %d in file with %d chunks", chunk, len(stco.ChunkOffsets)))
	}

	return stco.ChunkOffsets[chunk-1]
}
