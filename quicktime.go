package quicktime

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)


type Atom struct {
	Header AtomHeader
	Data []byte
}

func ParseAtomHeader(buffer []byte) (AtomHeader, error) {
	if len(buffer) < AtomHeaderLength {
		return AtomHeader{}, errors.New("Invalid buffer size")
	}

	// Read atom size
	var atomSize uint32
	if err := binary.Read(bytes.NewReader(buffer), binary.BigEndian, &atomSize); err != nil {
		return AtomHeader{}, err
	}

	if atomSize == 0 {
		return AtomHeader{}, errors.New("Zero size not supported yet")
	}

	if atomSize == 1 {
		return AtomHeader{}, errors.New("64 bit atom size not supported yet")
	}

	// Read atom type
	atomType := string(buffer[4:8])

	return AtomHeader{Offset: 0,
										Size: int(atomSize),
										DataSize: int(atomSize) - AtomHeaderLength,
										Type: atomType}, nil
}

func ParseAtomHeaderAt( r io.ReaderAt, offset int64 ) (AtomHeader, error) {
	header_buf := make([]byte, AtomHeaderLength )
	n,err := r.ReadAt(header_buf, offset)

	if err != nil {
		return AtomHeader{}, err
	}

	if n != AtomHeaderLength {
		return AtomHeader{}, errors.New("Couldn't read AtomHeaderLength")
	}

	atom,err := ParseAtomHeader( header_buf )
	atom.Offset = offset
	return atom,err
}

func ReadAtom(r io.Reader) (*Atom, error) {
	var atomBuffer bytes.Buffer

	// Read header
	if _, err := io.CopyN(&atomBuffer, r, AtomHeaderLength); err != nil {
		return nil, err
	}

	// Parse header
	atomHeader, err := ParseAtomHeader(atomBuffer.Bytes())
	if err != nil {
		return nil, err
	}

	// Read atom data
	if _, err = io.CopyN(&atomBuffer, r, int64(atomHeader.DataSize)); err != nil {
		return nil, err
	}

	// Create atom
	atom := Atom{atomHeader, atomBuffer.Bytes()}

	return &atom, nil
}

func ReadAtomAt( r io.ReaderAt, header AtomHeader ) (*Atom,error) {

	buf := make([]byte, header.DataSize )
	n,err := r.ReadAt( buf, header.Offset )
	if err != nil || n != header.DataSize {
		return nil, errors.New("Read incorrect number of bytes while getting Atom data")
	}

	return &Atom{ Header: header, Data: buf}, nil
}



// func main() {
//
// 	var init *IsoBmffInitSegment
// 	var err error
//
// 	for {
// 		if init == nil {
// 			init, err = ReadIsoBmffInitSegment(os.Stdin)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 		}
//
// 		log.Println(init.FTYP.Header, init.MOOV.Header)
//
// 		media, err := ReadIsoBmffMediaSegment(os.Stdin)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		log.Println(media.MOOF.Header, media.MDAT.Header)
// 	}
// }
