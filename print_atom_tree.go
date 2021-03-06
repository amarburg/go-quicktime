package quicktime

import "fmt"

// DumpTree iterates over an AtomArray and prints each atom on a line
func DumpTree(tree AtomArray) {
	fmt.Println(" (* = has loaded contents)")
	indent := 0
	for _, atom := range tree {
		PrintAtom(atom, indent)
	}
}

// PrintAtom writes a short descriptive string for each Atom in an AtomArray
// to do this, it must parse the atom
func PrintAtom(atom *Atom, indent int) {

	for i := 0; i < indent; i++ {
		fmt.Printf("..")
	}

	star := " "
	if atom.HasData() {
		star = "*"
	}
	fmt.Printf("%v%s  %v", atom.Type, star, atom.Size)

	// Atom-specific debug output
	if atom.HasData() {
		switch atom.Type {
		case "ftyp":
			ftyp, _ := ParseFTYP(atom)

			fmt.Printf(" (%08x %08x) ", ftyp.MajorBrand, ftyp.MinorVersion)
		case "stco", "co64":
			stco, _ := ParseSTCO(atom)

			fmt.Printf(" (%d entries) ", len(stco.ChunkOffsets))
		case "stsz":
			stsz, _ := ParseSTSZ(atom)

			fmt.Printf(" (%d entries, default size %d)", len(stsz.SampleSizes), stsz.SampleSize)
		case "stsc":
			stsz, _ := ParseSTSC(atom)

			fmt.Printf(" (%d entries): ", len(stsz.Entries))
			for _, e := range stsz.Entries {
				fmt.Printf(" (%d,%d,%d) ", e.FirstChunk, e.SamplesPerChunk, e.SampleID)
			}
		case "tkhd":
			tkhd, _ := ParseTKHD(atom)
			fmt.Printf(" (w: %.4f, h: %.4f)", tkhd.Width, tkhd.Height)
		case "mvhd":
			mvhd, _ := ParseMVHD(atom)
			seconds := float32(mvhd.TimeTicks) / float32(mvhd.TimeScale)
			fmt.Printf(" scale: %d, duration: %d (%.4f sec)", mvhd.TimeScale, mvhd.TimeTicks, seconds)
		}
	}

	fmt.Printf("\n")

	for _, child := range atom.Children {
		PrintAtom(child, indent+1)
	}

}
