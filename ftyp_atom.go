package quicktime

import "errors"

type FTYPAtom struct {
  Header AtomHeader
  MajorBrand uint32
  MinorVersion uint32
}


func ParseFTYP( atom *Atom ) (FTYPAtom,error) {
  if atom.Header.Type != "ftyp"{
    return FTYPAtom{}, errors.New("Not an FTYP atom")
  }

  return FTYPAtom{ Header: atom.Header,
                    MajorBrand: Uint32Decode( atom.Data[0:4] ),
                    MinorVersion: Uint32Decode( atom.Data[4:8] ) }, nil
}
