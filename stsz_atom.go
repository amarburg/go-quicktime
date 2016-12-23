package quicktime

import "bytes"
import "encoding/binary"
import "errors"
import "fmt"


type STSZAtom struct {
  Atom  *Atom
  SampleSze int32
  SampleSizes []int32
}

func ParseSTSZ( atom *Atom ) (STSZAtom, error) {
  if atom.Type != "stsz"{
    return STSZAtom{}, errors.New("Not an STSZ atom")
  }

  if !atom.HasData() {
    return STSZAtom{}, errors.New("STSZ Atom doesn't have data")
  }

  stsz := STSZAtom{ Atom: atom,
                    SampleSze: Int32Decode( atom.Data[4:8] ),
                    SampleSizes: make( []int32, 0 )}

  if( stsz.SampleSze == 0 ) {
    num_entries := Int32Decode( atom.Data[8:12] )

    stsz.SampleSizes = make([]int32, num_entries )

    if( num_entries > 0 ) {
      buf := bytes.NewBuffer( atom.Data[12:] )
      binary.Read( buf, binary.BigEndian, &stsz.SampleSizes )
    }
  }

  return stsz, nil
}

func (stsz STSZAtom) NumSamples() int {
  // How do you find length if only the universal sample size is provided?
  return len(stsz.SampleSizes)
}

func (stsz STSZAtom) SampleSize( sample int ) int {
  if stsz.SampleSze > 0 {return int(stsz.SampleSze)}

  if sample < 1 || sample > stsz.NumSamples() {
    panic(fmt.Sprintf("Requested sample %d in video with %d samples",sample, stsz.NumSamples() ) )
  }

  return int(stsz.SampleSizes[sample-1])
}
