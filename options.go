package worker

// Options holds options from command line
type Options map[string]string

// HasOneOf returns true if Options contains at least one of requested values
// Example:
// o.HasOneOf("q", "quiet")
// o.HasOneOf("v", "verbose")
func (this Options) HasOneOf(a ...string) bool {
	for k, _ := range this {
		for _, v := range a {
			if k == v {
				return true
			}
		}
	}

	return false
}
