package quicktime

import "io"

import "fmt"

type AtomArray []*Atom

// Functions for populating a tree of Atoms
// TODO:   Eager loading of some Atoms while building tree
func BuildTree(r io.ReaderAt, filesize int64) (AtomArray, error) {
	root := make([]*Atom, 0, 5)
	var err error = nil

	var offset int64 = 0
	for {
		atom, err := ParseAtomAt(r, offset)

		if err != nil {
			break
		}

		// Cheap eagerload...
		if atom.Type == "moov" {
			atom.LoadData(r)
		}

		if atom.IsContainer() {
			if atom.HasData() {
				atom.BuildChildren()
			} else {
				atom.ReadChildren(r)
			}
		}

		offset += int64(atom.Size)
		root = append(root, &atom)

	}
	return root, err
}

func (atom *Atom) ReadChildren(r io.ReaderAt) {
	var offset int64 = AtomHeaderLength
	for offset < int64(atom.Size) {
		loc := atom.Offset + offset
		//fmt.Println("Looking for header at:",loc)
		hdr, err := ParseAtomAt(r, loc)

		if err == nil {
			//fmt.Println("Found header at",loc,":", hdr.Type)
			if hdr.IsContainer() {
				hdr.ReadChildren(r)
			}

			offset += int64(hdr.Size)

			atom.Children = append(atom.Children, &hdr)

		} else {
			break
		}
	}
}

// As above, but uses parent.Data rather than the reader
func (atom *Atom) BuildChildren() {

	var offset int64 = 0
	for offset+AtomHeaderLength < int64(atom.Size) {
		fmt.Println("Looking for header at:", offset)
		hdr, err := ParseAtom(atom.Data[offset : offset+AtomHeaderLength])

		if err == nil {
			fmt.Println("Found header at", offset, ":", hdr.Type)
			hdr.Data = atom.Data[offset+AtomHeaderLength : offset+int64(hdr.Size)]

			if hdr.IsContainer() {
				hdr.BuildChildren()
			}

			offset += int64(hdr.Size)

			atom.Children = append(atom.Children, &hdr)

		} else {
			fmt.Println("Error parsing atom:", err.Error())
			break
		}
	}
}
