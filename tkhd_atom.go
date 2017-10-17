package quicktime

import "errors"

// A TKHDAtom ...
type TKHDAtom struct {
	Atom          *Atom
	Height, Width float32
}

// ParseTKHD converts a generic "tkhd" Atom to a TKHD atom
func ParseTKHD(atom *Atom) (TKHDAtom, error) {
	if atom.Type != "tkhd" {
		return TKHDAtom{}, errors.New("Not an TKHD atom")
	}

	if !atom.HasData() {
		return TKHDAtom{}, errors.New("TKHD Atom doesn't have data")
	}

	tkhd := TKHDAtom{Atom: atom,
		Width:  fixed32Decode(atom.Data[76:80]),
		Height: fixed32Decode(atom.Data[80:84])}

	return tkhd, nil
}
