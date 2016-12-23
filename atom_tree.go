package quicktime

import "io"


type AtomArray []Atom

// Functions for populating a tree of Atoms
// TODO:   Eager loading of some Atoms while building tree
func BuildTree( r io.ReaderAt, filesize int64 ) AtomArray {
  root := make( []Atom, 0, 5 )

  var offset int64 = 0
  for {
    header,err := ParseAtomAt( r, offset )

    if err != nil {
      break
    }

    if header.IsContainer() {
      header.AddChildren( r, true )
    }

    offset += int64( header.Size )
    root = append( root, header )

  }
  return root
}

func (atom* Atom) AddChildren( r io.ReaderAt, recursive bool ) {
  var offset int64 = AtomLength
  for offset < int64(atom.Size) {
    loc := atom.Offset+offset
    //fmt.Println("Looking for header at:",loc)
    hdr,err := ParseAtomAt( r, loc )

    if err == nil {
      //fmt.Println("Found header at",loc,":", hdr.Type)
      if recursive && hdr.IsContainer() {
       hdr.AddChildren(r, recursive)
      }

      offset += int64(hdr.Size)

      atom.Children = append( atom.Children, hdr )

    } else {
      break
    }
  }
}
