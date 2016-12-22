package quicktime

import "bytes"
import "encoding/binary"
import "errors"


type STSZAtom struct {
  Header  AtomHeader
  SampleSize int32
  SampleSizes []int32
}

func ParseSTSZ( atom *Atom ) (STSZAtom, error) {
  if atom.Header.Type != "stsz"{
    return STSZAtom{}, errors.New("Not an STSZ atom")
  }

  stsz := STSZAtom{ Header: atom.Header,
                    SampleSize: Int32Decode( atom.Data[4:8] ),
                    SampleSizes: make( []int32, 0 )}

  if( stsz.SampleSize == 0 ) {
    num_entries := Int32Decode( atom.Data[8:12] )

    stsz.SampleSizes = make([]int32, num_entries )

    if( num_entries > 0 ) {
      buf := bytes.NewBuffer( atom.Data[12:] )
      binary.Read( buf, binary.BigEndian, &stsz.SampleSizes )
    }
  }

  return stsz, nil
}
