package quicktime

import "io"
import "fmt"

func DumpTree( file io.ReaderAt, tree AtomArray ) {
  indent := 0
  for _,atom := range tree  {
    PrintAtom( file, atom, indent )
  }
}

func PrintAtom( file io.ReaderAt, atom *Atom, indent int ){

    for i := 0; i < indent; i++ { fmt.Printf("..") }
    fmt.Printf("%v  %v",atom.Type, atom.Size )

    switch atom.Type {
    case "ftyp":
      atom.LoadData( file )
      ftyp,_ := ParseFTYP( atom )

      fmt.Printf(" (%08x %08x) ", ftyp.MajorBrand, ftyp.MinorVersion )
    case "stco":
      atom.LoadData( file )
      stco,_ := ParseSTCO( atom )

      fmt.Printf(" (%d entries) ", len(stco.ChunkOffsets) )
    case "stsz":
      atom.LoadData( file )
      stsz,_ := ParseSTSZ( atom )

      fmt.Printf(" (%d entries, default size %d)", len(stsz.SampleSizes), stsz.SampleSize )
    case "stsc":
      atom.LoadData( file )
      stsz,_ := ParseSTSC( atom )

      fmt.Printf(" (%d entries): ", len(stsz.Entries) )
      for _,e := range stsz.Entries  {
        fmt.Printf(" (%d,%d,%d) ", e.FirstChunk, e.SamplesPerChunk, e.SampleId )
      }
    case "tkhd":
      atom.LoadData( file )
      tkhd,_ := ParseTKHD( atom )

      fmt.Printf(" (w: %.4f, h: %.4f)", tkhd.Width, tkhd.Height)
    }

    fmt.Printf("\n")

    for _,child := range atom.Children  {
      PrintAtom( file, child, indent+1 )
    }

}
