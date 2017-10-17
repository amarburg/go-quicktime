package quicktime

import "errors"

// A FTYPAtom stores versioning information about the Atom tree
type FTYPAtom struct {
	Atom         *Atom
	MajorBrand   uint32
	MinorVersion uint32
}

// ParseFTYP converts a generic "ftyp" atom to an FTYPAtom
func ParseFTYP(atom *Atom) (FTYPAtom, error) {
	if atom.Type != "ftyp" {
		return FTYPAtom{}, errors.New("Not an FTYP atom")
	}

	if !atom.HasData() {
		return FTYPAtom{}, errors.New("FTYP doesn't have data")
	}

	return FTYPAtom{Atom: atom,
		MajorBrand:   uint32Decode(atom.Data[0:4]),
		MinorVersion: uint32Decode(atom.Data[4:8])}, nil
}
