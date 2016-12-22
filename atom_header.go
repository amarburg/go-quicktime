package quicktime


const AtomHeaderLength = 8

type AtomHeader struct {
  Offset    int64
	Size     int
	DataSize int
	Type     string
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
