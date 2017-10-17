package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "io"
import "fmt"

// AtomHeaderLength is the length of a standard Atom header in bytes:
// a 4 byte atom size, and a 4 byte atom type.
const AtomHeaderLength = 8

// ExtendedHeaderLength is the length of an extended Quicktime Header:
// 4 bytes of size == 1 for extended headers
// 4 bytes of atom type
// 8 bytes of extended size
const ExtendedHeaderLength = 16

// AtomArray is used to store the topmost level when building the atom tree ... there is no master top-level Atom in Quicktime.
type AtomArray []*Atom

// The Atom struct stores a generic (not-decoded) Atom, including its type,
// size and optionally its Children and Data as a slice of the original byte buffer.
// If read from an io.ReaderAt it will also store the offset of the atom within the file
// (see ReadAtomAt)
type Atom struct {
	Offset   uint64
	Size     uint64
	DataSize uint64
	Type     string
	IsExt    bool

	Children AtomArray
	Data     []byte `json:"-"`
}

// HeaderLength returns the number of bytes in the header, depending on whether
// the header is extended or not.
func (atom *Atom) HeaderLength() uint64 {
	if atom.IsExt {
		return ExtendedHeaderLength
	}

	return AtomHeaderLength
}

// IsContainer returns true if a given atom type is a container, rather
// than a leaf.  This is done using a switch ... if I've missed an atom,
// send me a pull request!
func (atom *Atom) IsContainer() bool {
	switch atom.Type {
	case "moov", "trak", "mdia", "minf", "stbl", "dinf":
		return true
	}
	return false
}

// IsType returns true if the Atom's type matches the argument.
func (atom Atom) IsType(typeStr string) bool {
	return atom.Type == typeStr
}

// ParseAtom reads the first 8 bytes of the buffer and returns the appropriate Atom.
// Note this function doesn't set Data or Children for the Atom -- see ReadData and BuildChildren.
func ParseAtom(buffer []byte) (Atom, error) {
	if len(buffer) < AtomHeaderLength {
		return Atom{}, errors.New("Invalid buffer size")
	}

	atom := Atom{Offset: 0,
		IsExt: false,
		Type:  string(buffer[4:8])}

	// Read atom size
	var size32 uint32
	if err := binary.Read(bytes.NewReader(buffer), binary.BigEndian, &size32); err != nil {
		return Atom{}, err
	}

	atom.Size = uint64(size32)

	if atom.Size == 0 {
		return atom, errors.New("Zero size not supported yet")
	}

	if atom.Size == 1 {
		//bytes 8-15   atom size (including 16-byte size and type preamble)
		if len(buffer) < ExtendedHeaderLength {
			return Atom{}, errors.New("Got an extended-length header, but not enough data provided")
		}

		atom.IsExt = true

		var size64 uint64
		if err := binary.Read(bytes.NewReader(buffer[8:]), binary.BigEndian, &size64); err != nil {
			return atom, err
		}
		atom.Size = size64
	}

	atom.DataSize = atom.Size - atom.HeaderLength()

	return atom, nil
}

// ReadAtomAt reads the atom header from an io.ReaderAt and produces an Atom.
func ReadAtomAt(r io.ReaderAt, offset uint64) (Atom, error) {
	headerBuf := make([]byte, ExtendedHeaderLength)

	//fmt.Printf("Reading %d header bytes at offset %d\n", ExtendedHeaderLength, offset)
	//startTime := time.Now()
	n, err := r.ReadAt(headerBuf, int64(offset))

	//fmt.Printf("HTTP read took %d\n", time.Since(startTime) )

	if err != nil {
		fmt.Printf("Error reading atom at %d: %s\n", offset, err.Error())
		return Atom{}, err
	}

	if n != ExtendedHeaderLength {
		return Atom{}, fmt.Errorf("Couldn't read ExtendedHeaderLength bytes (%d) at offset %d, read %d bytes", ExtendedHeaderLength, offset, n)
	}

	atom, err := ParseAtom(headerBuf)
	atom.Offset = offset
	return atom, err
}

// ReadData populates the atom's Data member from the reader.   Assumes it is the same reader
// used to populate the original Atom.
// If Atom.Children is set, it also sets the Data for its children (recursively).
func (atom *Atom) ReadData(r io.ReaderAt) (err error) {
	if atom.HasData() {
		return nil
	}

	atom.Data = make([]byte, atom.DataSize)
	dataOffset := atom.Offset + atom.HeaderLength()

	//fmt.Printf("Reading %d data bytes at offset %d\n", atom.DataSize, dataOffset)
	//startTime := time.Now()
	n, err := r.ReadAt(atom.Data, int64(dataOffset))
	//fmt.Printf("HTTP read took %d\n", time.Since(startTime) )

	if err != nil && err.Error() != "EOF" {
		return err
	} else if uint64(n) != atom.DataSize {
		return fmt.Errorf("Read incorrect number of bytes while getting Atom data %d != %d", n, atom.DataSize)
	}

	for _, child := range atom.Children {
		child.SetData(atom.Data[child.Offset-dataOffset:])
	}

	return nil
}

// SetData sets the atom's Data field from the byte slice.   The slice must start
// at the beginning of this atom's header.
// If the Atom has children, it will set the data for the Children (recursively).
func (atom *Atom) SetData(buf []byte) {
	// TODO.   Check buf is atom.Size
	if uint64(len(buf)) < atom.Size {
		return
	}

	atom.Data = buf[atom.HeaderLength():atom.Size]

	for _, child := range atom.Children {
		child.SetData(buf[child.Offset-atom.Offset:])
	}
}

// HasData is true if the Data field is set for the atom
func (atom Atom) HasData() bool {
	return len(atom.Data) > 0
}
