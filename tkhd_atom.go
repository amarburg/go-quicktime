package quicktime

import "errors"



type TKHDAtom struct {
  Atom  *Atom
  Height, Width  float32
}

func ParseTKHD( atom *Atom ) (TKHDAtom, error) {
  if atom.Type != "tkhd"{
    return TKHDAtom{}, errors.New("Not an TKHD atom")
  }

  if !atom.HasData() {
    return TKHDAtom{}, errors.New("TKHD Atom doesn't have data")
  }

  tkhd := TKHDAtom{ Atom: atom,
                    Width: Fixed32Decode( atom.Data[76:80] ),
                    Height: Fixed32Decode( atom.Data[80:84] ) }


  return tkhd, nil
}
