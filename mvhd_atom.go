package quicktime

import "errors"

// A MVHDAtom stores timing information about the movie
type MVHDAtom struct {
	Atom                                        *Atom
	CreationTime, ModTime, TimeScale, TimeTicks uint32
}

// ParseMVHD converts a generic "mvhd" atom to an MVHDAtom
func ParseMVHD(atom *Atom) (MVHDAtom, error) {
	if atom.Type != "mvhd" {
		return MVHDAtom{}, errors.New("Not an MVHD atom")
	}

	if !atom.HasData() {
		return MVHDAtom{}, errors.New("MVHD Atom doesn't have data")
	}

	mvhd := MVHDAtom{Atom: atom,
		CreationTime: uint32Decode(atom.Data[4:8]),
		ModTime:      uint32Decode(atom.Data[8:12]),
		TimeScale:    uint32Decode(atom.Data[12:16]),
		TimeTicks:    uint32Decode(atom.Data[16:20]),
	}

	return mvhd, nil
}

// Duration returns the movie duration stored in an MVHDAtom
func (mvhd MVHDAtom) Duration() float32 {
	return float32(mvhd.TimeTicks) / float32(mvhd.TimeScale)
}
