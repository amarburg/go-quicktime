package quicktime

import "bytes"
import "encoding/binary"

func Uint32Decode( data []byte ) (ret uint32) {
  if len(data) < 4 { return 0; }

  buf := bytes.NewBuffer(data)
  binary.Read(buf, binary.BigEndian, &ret)
  return
}
