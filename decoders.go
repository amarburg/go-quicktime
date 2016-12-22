package quicktime

import "bytes"
import "encoding/binary"


func Int32Decode( data []byte ) (ret int32) {
  if len(data) < 4 { return 0; }

  buf := bytes.NewBuffer(data)
  binary.Read(buf, binary.BigEndian, &ret)
  return
}

func Uint32Decode( data []byte ) (ret uint32) {
  if len(data) < 4 { return 0; }

  buf := bytes.NewBuffer(data)
  binary.Read(buf, binary.BigEndian, &ret)
  return
}
