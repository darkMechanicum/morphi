package readers

// Abstract interface for anything that can
// get runes.
type RuneReader interface {

	// Get next rune, or I/O error. If no rune is available,
	// then nil is returned.
	NextRune() (*rune, error)

	// Add some runes back to reader.
	PushBack(runes string)
}
