package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "math"

type stscEntry struct {
	FirstChunk      int32
	SamplesPerChunk int32
	SampleID        int32
}

// A STSCAtom stores a table of samples per chunk
type STSCAtom struct {
	Atom    *Atom
	Entries []stscEntry
}

// ParseSTSC converts a generic "stsc" Atom to a STSCAtom
func ParseSTSC(atom *Atom) (STSCAtom, error) {
	if atom.Type != "stsc" {
		return STSCAtom{}, errors.New("Not an STSC atom")
	}

	if !atom.HasData() {
		return STSCAtom{}, errors.New("STSC atom doesn't have data")
	}

	numEntries := int32Decode(atom.Data[4:8])
	stsc := STSCAtom{Atom: atom,
		Entries: make([]stscEntry, numEntries)}

	if numEntries > 0 {
		buf := bytes.NewBuffer(atom.Data[8:])
		binary.Read(buf, binary.BigEndian, &stsc.Entries)
	}

	return stsc, nil
}

// SampleChunk converts a sample number to a chunk, chunk offset,
// and offset with the chunk
func (stsc STSCAtom) SampleChunk(s int) (chunk int, chunkStart int, offset int, err error) {

	// Convert to base zero
	s--

	// Hmm, what's an ideal way to handle this?
	switch len(stsc.Entries) {
	case 0:
		return -1, -1, -1, nil
	case 1:
		//  If there's one entry, then all chunks are the same size

		// (do all the math based zero)
		chunk = int(math.Floor(float64(s) / float64(stsc.Entries[0].SamplesPerChunk)))
		start := int(chunk * int(stsc.Entries[0].SamplesPerChunk))
		offset = s - start

		// Convert back to base one
		chunk++
		offset++
		start++

		//fmt.Printf("STSC: I believe frame %d is sample %d in chunk %d (which starts at %d)\n", s+1, offset, chunk, start)

		// Remember that chunks and samples are 1-based
		return chunk, start, offset, nil
	}

	sample := int32(s)

	lastChunk := stsc.Entries[len(stsc.Entries)-1].FirstChunk
	samplesPerChunk := stsc.Entries[0].SamplesPerChunk
	nextEntry := 0
	accum := int32(0)

	for i := int32(0); i <= lastChunk; i++ {
		if i == stsc.Entries[nextEntry].FirstChunk {
			samplesPerChunk = stsc.Entries[nextEntry].SamplesPerChunk
			nextEntry++
		}

		if samplesPerChunk > (sample) {
			return int(i + 1), int(accum + 1), int(sample + 1), nil
		}

		sample -= samplesPerChunk
		accum += samplesPerChunk
	}

	// If you get here, you're on the last chunk
	chunk = int(lastChunk) + int(math.Floor(float64(sample)/float64(samplesPerChunk)))
	offset = int(math.Remainder(float64(sample), float64(samplesPerChunk)))

	// Convert back to base one
	chunk++
	offset++
	accum++

	//fmt.Printf("STSC: I believe frame %d is sample %d in chunk %d (which starts at %d)\n", s+1, offset, chunk, accum)

	return chunk, int(accum), offset, nil
}
