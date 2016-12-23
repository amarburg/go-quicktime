package quicktime

import "errors"

type STBLAtom struct {
  stsz  STSZAtom
  stco  CO64Atom
  stsc  STSCAtom
}


func ParseSTBL( atom Atom ) (STBLAtom, error) {
  if atom.Type != "stbl"{
    return STBLAtom{}, errors.New("Not an STBL atom")
  }

  if atom.HasData() {
    return STBLAtom{}, errors.New( "STBL Atom doesn't have data" )
  }

  stbl := STBLAtom{}

  // Having data implies all children have data as well
  for _,child := range atom.Children {
    switch child.Type {
    case "stsc": stbl.stsc,_ = ParseSTSC( child )
    case "stco": stbl.stco,_ = ParseSTCO( child )
    case "stsz": stbl.stsz,_ = ParseSTSZ( child )
    }
  }

  return stbl, nil
}
