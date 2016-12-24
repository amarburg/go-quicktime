package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "io"

const AtomHeaderLength = 8

type Atom struct {
	Offset   int64
	Size     int
	DataSize int
	Type     string

	Children []*Atom
	Data     []byte
}

func (header *Atom) IsContainer() bool {
	switch header.Type {
	case "moov", "trak", "mdia", "minf", "stbl", "dinf":
		return true
	}
	return false
}

func (header Atom) IsType(type_str string) bool {
	return header.Type == type_str
}

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
		return Atom{}, errors.New("64 bit atom size not supported yet")
	}

	// Read atom type
	atomType := string(buffer[4:8])

	return Atom{Offset: 0,
		Size:     int(atomSize),
		DataSize: int(atomSize) - AtomHeaderLength,
		Type:     atomType}, nil
}

func ParseAtomAt(r io.ReaderAt, offset int64) (Atom, error) {
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

// TODO:   Loading data should also populate Data in children
func (atom *Atom) LoadData(r io.ReaderAt) (err error) {
	if atom.HasData() {
		return nil
	}

	buf := make([]byte, atom.DataSize)
	n, err := r.ReadAt(buf, atom.Offset+AtomHeaderLength)
	if err != nil || n != atom.DataSize {
		return errors.New("Read incorrect number of bytes while getting Atom data")
	}
	atom.Data = buf
	return nil
}

func (atom Atom) HasData() bool {
	return len(atom.Data) > 0
}
