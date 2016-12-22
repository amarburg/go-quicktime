package quicktime

import "bytes"
import "encoding/binary"
import "errors"

type STSZAtom struct {
  Header  AtomHeader
  SampleSize uint32
  SampleSizes []int32
}

func ParseSTSZ( atom *Atom ) (STSZAtom, error) {
  if atom.Header.Type != "stsz"{
    return STSZAtom{}, errors.New("Not an STSZ atom")
  }

  num_entries := Uint32Decode( atom.Data[8:12] )
  stsz := STSZAtom{ Header: atom.Header,
                    SampleSize: Uint32Decode( atom.Data[4:8] ),
                    SampleSizes: make([]int32, num_entries ) }

  buf := bytes.NewBuffer( atom.Data[12:] )
  binary.Read( buf, binary.BigEndian, &stsz.SampleSizes )

  return stsz, nil
}
