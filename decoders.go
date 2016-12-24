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

func Fixed32Decode( data []byte ) (ret float32) {
  if len(data) < 4 { return 0; }

  ints := [2]int16{0,0}
  buf := bytes.NewBuffer(data)
  binary.Read(buf, binary.BigEndian, &ints)

  return float32(ints[0]) + float32(ints[1])/65536.0
}
