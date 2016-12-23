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
  Atom     Atom
  Entries []STSCEntry
}

func ParseSTSC( atom Atom ) (STSCAtom, error) {
  if atom.Type != "stsc"{
    return STSCAtom{}, errors.New("Not an STSC atom")
  }

  if !atom.HasData() {
    return STSCAtom{}, errors.New("STSC atom doesn't have data")
  }

  num_entries := Int32Decode( atom.Data[4:8])
  stsc := STSCAtom{ Atom: atom,
                    Entries: make( []STSCEntry, num_entries ) }

  if( num_entries > 0 ) {
    buf := bytes.NewBuffer( atom.Data[8:] )
    binary.Read( buf, binary.BigEndian, &stsc.Entries )
  }


  return stsc, nil
}
