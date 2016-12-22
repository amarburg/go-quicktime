package quicktime

import "bytes"
import "encoding/binary"
import "errors"


type STCOAtom struct {
  Header  AtomHeader
  ChunkOffsets []int64
}

func ParseSTCO( atom *Atom ) (STCOAtom, error) {
  if atom.Header.Type != "stco"{
    return STCOAtom{}, errors.New("Not an STCO atom")
  }

  num_entries := Uint32Decode( atom.Data[4:8] )

  stco := STCOAtom{ Header: atom.Header,
                    ChunkOffsets: make([]int64, num_entries ) }
  buf := bytes.NewBuffer( atom.Data[8:] )
  binary.Read( buf, binary.BigEndian, &stco.ChunkOffsets)

  return stco, nil
}
