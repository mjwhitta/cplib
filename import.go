package cplib

// Import is a generic struct that contains the imported function name
// and the library that is expected to contain it.
type Import struct {
	Lib  string
	Name string
}
