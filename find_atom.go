package quicktime

func (atom Atom) FindAtom(type_str string) *Atom {
	//fmt.Println(atom.Type,"atom has", len(atom.Children),"children when looking for", type_str)
	for _, child := range atom.Children {
		//fmt.Println("Comparing",atom.Type,"to",type_str)
		if child.IsType(type_str) {
			return child
		}
	}
	return nil
}

func (atom Atom) FindAtoms(type_str string) (out []*Atom) {
	out = make([]*Atom, 0)

	for _, child := range atom.Children {
		if child.IsType(type_str) {
			//fmt.Println("Found atom",child.Type,"we were looking for at", child)
			out = append(out, child)
		}
	}
	return out
}

func (atoms AtomArray) FindAtom(type_str string) *Atom {
	for _, child := range atoms {
		if child.IsType(type_str) {
			return child
		}
	}
	return nil
}

func (atoms AtomArray) FindAtoms(type_str string) []*Atom {
	out := make([]*Atom, 0)

	for _, child := range atoms {
		if child.IsType(type_str) {
			out = append(out, child)
		}
	}
	return out
}
