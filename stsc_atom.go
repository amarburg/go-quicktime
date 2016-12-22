package quicktime

import "bytes"
import "encoding/binary"
import "errors"

type STSCEntry struct {
  FirstChunk int32
  SamplesPerChunk int32
  SampleId int32
}

type STSCAtom struct {
  Header     AtomHeader
  Entries []STSCEntry
}

func ParseSTSC( atom *Atom ) (STSCAtom, error) {
  if atom.Header.Type != "stsc"{
    return STSCAtom{}, errors.New("Not an STSC atom")
  }

  num_entries := Int32Decode( atom.Data[4:8])
  stsc := STSCAtom{ Header: atom.Header,
                    Entries: make( []STSCEntry, num_entries ) }

  if( num_entries > 0 ) {
    buf := bytes.NewBuffer( atom.Data[8:] )
    binary.Read( buf, binary.BigEndian, &stsc.Entries )
  }


  return stsc, nil
}
