package quicktime

import "bytes"
import "encoding/binary"
import "errors"

type FTYPAtom struct {
  Header AtomHeader
  MajorBrand uint32
  MinorVersion uint32
}


func Uint32Decode( data []byte ) (ret uint32) {
  if len(data) < 4 { return 0; }

  buf := bytes.NewBuffer(data)
  binary.Read(buf, binary.BigEndian, &ret)
  return
}

func ParseFTYP( atom *Atom ) (FTYPAtom,error) {
  if atom.Header.Type != "ftyp"{
    return FTYPAtom{}, errors.New("Not an FTYP atom")
  }

  return FTYPAtom{ Header: atom.Header,
                    MajorBrand: Uint32Decode( atom.Data[0:4] ),
                    MinorVersion: Uint32Decode( atom.Data[4:8] ) }, nil
}
