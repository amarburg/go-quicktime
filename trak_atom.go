package quicktime

import "errors"

type MINFAtom struct {
	Stbl STBLAtom
}

type MDIAAtom struct {
	Minf MINFAtom
}

type TRAKAtom struct {
	Tkhd TKHDAtom
	Mdia MDIAAtom
}

func ParseTRAK(atom *Atom) (TRAKAtom, error) {
	if atom.Type != "trak" {
		return TRAKAtom{}, errors.New("Not an TRAK atom")
	}

	stbl_atom := atom.FindAtom("mdia").FindAtom("minf").FindAtom("stbl")

	// Skip the intermediate layers for now
	stbl, err := ParseSTBL(stbl_atom)
	if err != nil {
		return TRAKAtom{}, err
	}

	tkhd, err := ParseTKHD(atom.FindAtom("tkhd"))
	if err != nil {
		return TRAKAtom{}, err
	}

	trak := TRAKAtom{Tkhd: tkhd,
		Mdia: MDIAAtom{Minf: MINFAtom{Stbl: stbl}}}

	return trak, nil
}
