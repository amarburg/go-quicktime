// Package quicktime provides primitives for loading the tree of Atoms from a Quicktime file, and reading a subset of the atoms.
package quicktime

// import (
// 	"bytes"
// 	"errors"
// 	"io"
// )
//
//
// type Atom struct {
// 	Header AtomHeader
// 	Data []byte
// }
//
// func ReadAtom(r io.Reader) (*Atom, error) {
// 	var atomBuffer bytes.Buffer
//
// 	// Read header
// 	if _, err := io.CopyN(&atomBuffer, r, AtomHeaderLength); err != nil {
// 		return nil, err
// 	}
//
// 	// Parse header
// 	atomHeader, err := ParseAtomHeader(atomBuffer.Bytes())
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// Read atom data
// 	if _, err = io.CopyN(&atomBuffer, r, int64(atomHeader.DataSize)); err != nil {
// 		return nil, err
// 	}
//
// 	// Create atom
// 	atom := Atom{atomHeader, atomBuffer.Bytes()}
//
// 	return &atom, nil
// }
//
// func ReadAtomAt( r io.ReaderAt, header AtomHeader ) (*Atom,error) {
//
// 	buf := make([]byte, header.DataSize )
// 	n,err := r.ReadAt( buf, header.Offset + AtomHeaderLength )
// 	if err != nil || n != header.DataSize {
// 		return nil, errors.New("Read incorrect number of bytes while getting Atom data")
// 	}
//
// 	return &Atom{ Header: header, Data: buf}, nil
// }

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
