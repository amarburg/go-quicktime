package quicktime

import "bytes"
import "encoding/binary"
import "errors"


// I just cast everything to a CO64 atom for in-memory representation
type CO64Atom struct {
  Atom        Atom
  ChunkOffsets []int64
}


// Also handle CO64
func ParseSTCO( atom Atom ) (CO64Atom, error) {

  if !atom.HasData() {
    return CO64Atom{}, errors.New("STCO Atom doesn't have data")
  }

  // TODO:  Fix the DRY
  if atom.Type == "stco" {

    num_entries := Uint32Decode( atom.Data[4:8] )

    int32Offsets := make( []int32, num_entries )
    buf := bytes.NewBuffer( atom.Data[8:] )
    binary.Read( buf, binary.BigEndian, &int32Offsets)

    stco := CO64Atom{ Atom: atom,
                      ChunkOffsets: make([]int64, num_entries ) }

    // Copy the int32s to int64s
    // TODO:  I'm sure this isn't idiomatic
    for i :=0; i < int(num_entries); i++ {
      stco.ChunkOffsets[i] = int64( int32Offsets[i] )
    }

    return stco, nil


  } else if atom.Type == "co64" {


    num_entries := Uint32Decode( atom.Data[4:8] )

    stco := CO64Atom{ Atom: atom,
                      ChunkOffsets: make([]int64, num_entries ) }

    buf := bytes.NewBuffer( atom.Data[8:] )
    binary.Read( buf, binary.BigEndian, &stco.ChunkOffsets)

    return stco, nil

  } else {
    return CO64Atom{}, errors.New("Not an STCO atom")
  }


}
