package quicktime


func (atom Atom) FindAtom( type_str string ) (*Atom) {
  for _,child := range atom.Children {
    if child.IsType( type_str ) {
      return &child
    }
  }
  return nil
}

func (atom Atom) FindAtoms( type_str string ) ([]*Atom) {
  out := make([]*Atom,0)

  for _,child := range atom.Children {
    if child.IsType( type_str ) {
      out=append( out, &child )
    }
  }
  return out
}


func (atoms AtomArray) FindAtom( type_str string ) (*Atom) {
  for _,child := range atoms {
    if child.IsType( type_str ) {
      return &child
    }
  }
  return nil
}

func (atoms AtomArray) FindAtoms( type_str string ) ([]*Atom) {
  out := make([]*Atom,0)

  for _,child := range atoms {
    if child.IsType( type_str ) {
      out=append( out, &child )
    }
  }
  return out
}
