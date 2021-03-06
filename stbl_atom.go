package quicktime

import "errors"
import "fmt"

// A STBLAtom stores the three Atoms requires to look up the file offset of
// an individual frame ("sample") in the file:  STSZ (sample size),
// STCO or CO64 (chunk offset) and STSC (samples per chunk)
type STBLAtom struct {
	Stsz STSZAtom
	Stco CO64Atom
	Stsc STSCAtom
}

// ParseSTBL converts a generic Atom containing a "stbl" to a STBLAtom.
func ParseSTBL(atom *Atom) (STBLAtom, error) {
	if atom.Type != "stbl" {
		return STBLAtom{}, errors.New("Not an STBL atom")
	}

	if !atom.HasData() {
		return STBLAtom{}, errors.New("STBL Atom doesn't have data")
	}

	stbl := STBLAtom{}

	// Having data implies all children have data as well
	for _, child := range atom.Children {
		switch child.Type {
		case "stsc":
			stbl.Stsc, _ = ParseSTSC(child)
		case "stco", "co64":
			stbl.Stco, _ = ParseSTCO(child)
		case "stsz":
			stbl.Stsz, _ = ParseSTSZ(child)
		}
	}

	return stbl, nil
}

// SampleOffset calculates the byte offset of a given sample within the
// quicktime file
func (stbl STBLAtom) SampleOffset(sample int) (int64, error) {

	if uint64(sample) > stbl.NumFrames() {
		return 0, fmt.Errorf("Requested sample %d in file %d samples long", sample, stbl.NumFrames())
	}

	// Use STCO to determine which chunk it's in
	chunk, chunkStart, _, err := stbl.Stsc.SampleChunk(sample)
	if chunk < 0 {
		panic(fmt.Sprintln("Couldn't determine which chunk", sample, "is in"))
	} else if err != nil {
		panic(fmt.Sprintln("Error determining chunk: %s", err.Error()))
	}

	//fmt.Printf("STBL: Believe sample %d is number %d in chunk %d which starts with sample %d\n", sample, remainder, chunk, chunkStart)

	offset := stbl.Stco.ChunkOffset(chunk)
	for i := chunkStart; i < sample; i++ {
		offset += int64(stbl.Stsz.SampleSize(i))
	}

	return offset, nil
}

// NumFrames returns the number of samples (frames) in the stbl
func (stbl STBLAtom) NumFrames() uint64 {
	// What if there's only one sample?
	return stbl.Stsz.NumSamples()
}

// SampleOffsetSize calculates both the offset and size of a sample
func (stbl STBLAtom) SampleOffsetSize(sample int) (int64, int, error) {

	offset, err := stbl.SampleOffset(sample)
	sz := stbl.Stsz.SampleSize(sample)

	return offset, sz, err
}
