package quicktime

// StringList stores a slice of Strings for BuildTreeConfig
type StringList []string

// Tests if a string occurs in a StringList.
func (list StringList) includes(val string) bool {
	for _, str := range list {
		if str == val {
			return true
		}
	}
	return false
}
