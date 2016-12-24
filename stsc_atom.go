package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "math"

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
func (stsc STSCAtom) SampleChunk(s int) (int, int, int) {
	// Hmm, what's an ideal way to handle this?
	switch len(stsc.Entries) {
	case 0:
		return -1, -1, -1
	case 1:
		return 1, 1, s
	}

	sample := int32(s)

	last_chunk := stsc.Entries[len(stsc.Entries)-1].FirstChunk
	samples_per_chunk := stsc.Entries[0].SamplesPerChunk
	next_entry := 1
	accum := int32(1)

	for i := int32(1); i <= last_chunk; i++ {
		if i == stsc.Entries[next_entry].FirstChunk {
			samples_per_chunk = stsc.Entries[next_entry].SamplesPerChunk
			next_entry++
		}

		if samples_per_chunk > (sample - 1) {
			return int(i), int(accum), int(sample)
		}

		sample -= samples_per_chunk
		accum += samples_per_chunk
	}

	// If you get here, you're on the last chunk
	chunk := int(last_chunk) + int(math.Floor(float64(sample)/float64(samples_per_chunk)))
	offset := int(math.Remainder(float64(sample), float64(samples_per_chunk))) + 1
	return chunk, int(accum), offset
}
