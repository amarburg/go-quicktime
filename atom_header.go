package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "io"

const AtomHeaderLength = 8

type AtomHeader struct {
  Offset    int64
	Size     int
	DataSize int
	Type     string

  Children []AtomHeader
}

func (header* AtomHeader) IsContainer() bool {
	switch header.Type {
	case "moov","trak","mdia","minf","stbl","dinf": return true;
	}
	return false;
}

func (header AtomHeader) IsType( type_str string ) bool {
  return header.Type == type_str
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


func (atom* AtomHeader) AddChildren( r io.ReaderAt, recursive bool ) {
  var offset int64 = AtomHeaderLength
  for offset < int64(atom.Size) {
    loc := atom.Offset+offset
    //fmt.Println("Looking for header at:",loc)
    hdr,err := ParseAtomHeaderAt( r, loc )

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

func BuildTree( r io.ReaderAt, filesize int64 ) []AtomHeader {

  root := make( []AtomHeader, 0, 5 )

  var offset int64 = 0
  for {
    header,err := ParseAtomHeaderAt( r, offset )

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
