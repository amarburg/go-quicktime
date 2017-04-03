package quicktime

import "errors"

type MVHDAtom struct {
	Atom                                        *Atom
	CreationTime, ModTime, TimeScale, TimeTicks uint32
}

func ParseMVHD(atom *Atom) (MVHDAtom, error) {
	if atom.Type != "mvhd" {
		return MVHDAtom{}, errors.New("Not an MVHD atom")
	}

	if !atom.HasData() {
		return MVHDAtom{}, errors.New("MVHD Atom doesn't have data")
	}

	mvhd := MVHDAtom{Atom: atom,
		CreationTime: Uint32Decode(atom.Data[4:8]),
		ModTime:      Uint32Decode(atom.Data[8:12]),
		TimeScale:    Uint32Decode(atom.Data[12:16]),
		TimeTicks:    Uint32Decode(atom.Data[16:20]),
	}

	return mvhd, nil
}

func (mvhd MVHDAtom) Duration() float32 {
	return float32(mvhd.TimeTicks) / float32(mvhd.TimeScale)
}
