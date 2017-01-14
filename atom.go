package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "io"

// AtomHeaderLength is the length of a standard Atom header in bytes: a 4 byte atom size, and a 4 byte atom type.
const AtomHeaderLength = 8

// The Atom struct stores a generic (not-decoded) Atom, including its type, 
// size and optionally its Children and Data as a slice of the original byte buffer.
// If read from an io.ReaderAt it will also store the offset of the atom within the file 
// (see ReadAtomAt)
type Atom struct {
	Offset   int64
	Size     int
	DataSize int
	Type     string

	Children []*Atom
	Data     []byte
}

// IsContainer returns true if a given atom type is a container, rather than a leaf.  This is done using a switch ... if I've missed an atom, send me a pull request!
func (header *Atom) IsContainer() bool {
	switch header.Type {
	case "moov", "trak", "mdia", "minf", "stbl", "dinf":
		return true
	}
	return false
}

// IsType returns true if the Atom's type matches the argument.
func (header Atom) IsType(type_str string) bool {
	return header.Type == type_str
}

// ParseAtom reads the first 8 bytes of the buffer and returns the appropriate Atom.   
// Note this function doesn't set Data or Children for the Atom -- see ReadData and BuildChildren.
func ParseAtom(buffer []byte) (Atom, error) {
	if len(buffer) < AtomHeaderLength {
		return Atom{}, errors.New("Invalid buffer size")
	}

	// Read atom size
	var atomSize uint32
	if err := binary.Read(bytes.NewReader(buffer), binary.BigEndian, &atomSize); err != nil {
		return Atom{}, err
	}

	if atomSize == 0 {
		return Atom{}, errors.New("Zero size not supported yet")
	}

	if atomSize == 1 {
		//bytes 8-15   atom size (including 16-byte size and type preamble)
		return Atom{}, errors.New("64 bit atom size not supported yet")
	}

	// Read atom type
	atomType := string(buffer[4:8])

	return Atom{Offset: 0,
		Size:     int(atomSize),
		DataSize: int(atomSize) - AtomHeaderLength,
		Type:     atomType}, nil
}

// ReadAtom reads the atom header from an io.ReaderAt and produces an Atom.
func ReadAtomAt(r io.ReaderAt, offset int64) (Atom, error) {
	header_buf := make([]byte, AtomHeaderLength)
	n, err := r.ReadAt(header_buf, offset)

	if err != nil {
		return Atom{}, err
	}

	if n != AtomHeaderLength {
		return Atom{}, errors.New("Couldn't read AtomHeaderLength")
	}

	atom, err := ParseAtom(header_buf)
	atom.Offset = offset
	return atom, err
}

// ReadData populates the Atom's Data member from the reader.   Assumes it is the same reader
// used to populated the original Atom.
// If Atom.Children is set, it also sets the Data for its children (recursively).
func (atom *Atom) ReadData(r io.ReaderAt) (err error) {
	if atom.HasData() {
		return nil
	}

	buf := make([]byte, atom.DataSize)
	n, err := r.ReadAt(buf, atom.Offset+AtomHeaderLength)
	if err != nil || n != atom.DataSize {
		return errors.New("Read incorrect number of bytes while getting Atom data")
	}
	atom.Data = buf

	for _,child := range atom.Children {
		child.SetData( buf[ child.Offset - atom.Offset - AtomHeaderLength:] )
	}

	return nil
}

// SetData sets the atom's Data field from the byte slice.   The slice must start
// at the beginning of this atom's header.
// If the Atom has children, it will set the data for the Children (recursively).
func (atom *Atom) SetData(buf []byte) {
	// TODO.   Check buf is atom.Size
	if( len(buf) < atom.Size ) { return }
	
	atom.Data = buf[ AtomHeaderLength:atom.Size ]

	for _,child := range atom.Children {
		child.SetData( buf[ child.Offset - atom.Offset:] )
	}
}

// HasData is true if the Data field is set for the atom
func (atom Atom) HasData() bool {
	return len(atom.Data) > 0
}
