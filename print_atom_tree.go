package quicktime

import "fmt"

func DumpTree(tree AtomArray) {
	indent := 0
	for _, atom := range tree {
		PrintAtom(atom, indent)
	}
}

func PrintAtom(atom *Atom, indent int) {

	for i := 0; i < indent; i++ {
		fmt.Printf("..")
	}
	fmt.Printf("%v  %v", atom.Type, atom.Size)

	if atom.HasData() {
		switch atom.Type {
		case "ftyp":
			ftyp, _ := ParseFTYP(atom)

			fmt.Printf(" (%08x %08x) ", ftyp.MajorBrand, ftyp.MinorVersion)
		case "stco":
			stco, _ := ParseSTCO(atom)

			fmt.Printf(" (%d entries) ", len(stco.ChunkOffsets))
		case "stsz":
			stsz, _ := ParseSTSZ(atom)

			fmt.Printf(" (%d entries, default size %d)", len(stsz.SampleSizes), stsz.SampleSize)
		case "stsc":
			stsz, _ := ParseSTSC(atom)

			fmt.Printf(" (%d entries): ", len(stsz.Entries))
			for _, e := range stsz.Entries {
				fmt.Printf(" (%d,%d,%d) ", e.FirstChunk, e.SamplesPerChunk, e.SampleId)
			}
		case "tkhd":
			tkhd, _ := ParseTKHD(atom)

			fmt.Printf(" (w: %.4f, h: %.4f)", tkhd.Width, tkhd.Height)
		}
	}

	fmt.Printf("\n")

	for _, child := range atom.Children {
		PrintAtom(child, indent+1)
	}

}
