package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "fmt"

// A STSZAtom stores a table of sample sizes
type STSZAtom struct {
	Atom        *Atom
	SampleSze   int32
	SampleSizes []int32
}

// ParseSTSZ converts a generic "stsz" Atom to a STSZAtom
func ParseSTSZ(atom *Atom) (STSZAtom, error) {
	if atom.Type != "stsz" {
		return STSZAtom{}, errors.New("Not an STSZ atom")
	}

	if !atom.HasData() {
		return STSZAtom{}, errors.New("STSZ Atom doesn't have data")
	}

	stsz := STSZAtom{Atom: atom,
		SampleSze:   int32Decode(atom.Data[4:8]),
		SampleSizes: make([]int32, 0)}

	if stsz.SampleSze == 0 {
		numEntries := int32Decode(atom.Data[8:12])

		stsz.SampleSizes = make([]int32, numEntries)

		if numEntries > 0 {
			buf := bytes.NewBuffer(atom.Data[12:])
			binary.Read(buf, binary.BigEndian, &stsz.SampleSizes)
		}
	}

	return stsz, nil
}

// NumSamples returns the number of samples in an STSZAtom
func (stsz STSZAtom) NumSamples() int {
	// How do you find length if only the universal sample size is provided?
	return len(stsz.SampleSizes)
}

// SampleSize looks up the size of a given sample in the STSZAtom
func (stsz STSZAtom) SampleSize(sample int) int {
	if stsz.SampleSze > 0 {
		return int(stsz.SampleSze)
	}

	if sample < 1 || sample > stsz.NumSamples() {
		// TODO.   Duh, get rid of this panix and replace it with an error message...
		panic(fmt.Sprintf("Requested sample %d in video with %d samples", sample, stsz.NumSamples()))
	}

	return int(stsz.SampleSizes[sample-1])
}
