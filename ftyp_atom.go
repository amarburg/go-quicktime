package quicktime

import "errors"

type FTYPAtom struct {
  Atom      Atom
  MajorBrand uint32
  MinorVersion uint32
}


func ParseFTYP( atom Atom ) (FTYPAtom,error) {
  if atom.Type != "ftyp"{
    return FTYPAtom{}, errors.New("Not an FTYP atom")
  }

  if !atom.HasData() {
    return FTYPAtom{}, errors.New("FTYP doesn't have data")
  }

  return FTYPAtom{ Atom: atom,
                    MajorBrand: Uint32Decode( atom.Data[0:4] ),
                    MinorVersion: Uint32Decode( atom.Data[4:8] ) }, nil
}
