package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "math"

//import "fmt"

type STSCEntry struct {
	FirstChunk      int32
	SamplesPerChunk int32
	SampleId        int32
}

type STSCAtom struct {
	Atom    *Atom
	Entries []STSCEntry
}

func ParseSTSC(atom *Atom) (STSCAtom, error) {
	if atom.Type != "stsc" {
		return STSCAtom{}, errors.New("Not an STSC atom")
	}

	if !atom.HasData() {
		return STSCAtom{}, errors.New("STSC atom doesn't have data")
	}

	num_entries := Int32Decode(atom.Data[4:8])
	stsc := STSCAtom{Atom: atom,
		Entries: make([]STSCEntry, num_entries)}

	if num_entries > 0 {
		buf := bytes.NewBuffer(atom.Data[8:])
		binary.Read(buf, binary.BigEndian, &stsc.Entries)
	}

	return stsc, nil
}

// Remember that chunks and samples are 1-based
func (stsc STSCAtom) SampleChunk(s int) (chunk int, chunk_start int, offset int, err error) {

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

	last_chunk := stsc.Entries[len(stsc.Entries)-1].FirstChunk
	samples_per_chunk := stsc.Entries[0].SamplesPerChunk
	next_entry := 0
	accum := int32(0)

	for i := int32(0); i <= last_chunk; i++ {
		if i == stsc.Entries[next_entry].FirstChunk {
			samples_per_chunk = stsc.Entries[next_entry].SamplesPerChunk
			next_entry++
		}

		if samples_per_chunk > (sample) {
			return int(i + 1), int(accum + 1), int(sample + 1), nil
		}

		sample -= samples_per_chunk
		accum += samples_per_chunk
	}

	// If you get here, you're on the last chunk
	chunk = int(last_chunk) + int(math.Floor(float64(sample)/float64(samples_per_chunk)))
	offset = int(math.Remainder(float64(sample), float64(samples_per_chunk)))

	// Convert back to base one
	chunk++
	offset++
	accum++

	//fmt.Printf("STSC: I believe frame %d is sample %d in chunk %d (which starts at %d)\n", s+1, offset, chunk, accum)

	return chunk, int(accum), offset, nil
}
