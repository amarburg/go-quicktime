package quicktime

// FindAtom looks through an Atom's children for a given type,
// returning the first instance or nil of not found
func (atom Atom) FindAtom(typeStr string) *Atom {
	return atom.Children.FindAtom(typeStr)
}

// FindAtoms looks through an Atom's children for a given type,
// returning an array of all instances.  The array will be empty
// if not found.
func (atom Atom) FindAtoms(typeStr string) (out []*Atom) {
	return atom.Children.FindAtoms(typeStr)
}

// FindAtom looks through an AtomArray for a given type,
// returning the first instance or nil of not found
func (atoms AtomArray) FindAtom(typeStr string) *Atom {
	for _, child := range atoms {
		if child.IsType(typeStr) {
			return child
		}
	}
	return nil
}

// FindAtoms looks through an AtomArray for a given type,
// returning an array of all instances.  The array will be empty
// if not found.
func (atoms AtomArray) FindAtoms(typeStr string) AtomArray {
	var out AtomArray

	for _, child := range atoms {
		if child.IsType(typeStr) {
			out = append(out, child)
		}
	}
	return out
}
