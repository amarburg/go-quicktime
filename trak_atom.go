package quicktime

import "errors"

// A MINFAtom is ...
type MINFAtom struct {
	Stbl STBLAtom
}

// A MDIAAtom ...
type MDIAAtom struct {
	Minf MINFAtom
}

// A TRAKAtom ...
type TRAKAtom struct {
	Tkhd TKHDAtom
	Mdia MDIAAtom
}

// ParseTRAK converts a generic "trak" Atom to a TRAKAtom
func ParseTRAK(atom *Atom) (TRAKAtom, error) {
	if atom.Type != "trak" {
		return TRAKAtom{}, errors.New("Not an TRAK atom")
	}

	stblAtom := atom.FindAtom("mdia").FindAtom("minf").FindAtom("stbl")

	// Skip the intermediate layers for now
	stbl, err := ParseSTBL(stblAtom)
	if err != nil {
		return TRAKAtom{}, err
	}

	tkhd, err := ParseTKHD(atom.FindAtom("tkhd"))
	if err != nil {
		return TRAKAtom{}, err
	}

	trak := TRAKAtom{
		Tkhd: tkhd,
		Mdia: MDIAAtom{
			Minf: MINFAtom{
				Stbl: stbl,
			},
		},
	}

	return trak, nil
}
