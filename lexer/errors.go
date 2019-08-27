package lexer

type LexerError error

type NoRunesError struct{}

func (self *NoRunesError) Error() string {
	return "No runes are available for lexical parsing."
}

// Error, raising when buffer contains delimiter and non delimiter runes.
type DirtyBufferError struct{}

func (self *DirtyBufferError) Error() string {
	return "Buffer contains delimiter and non delimiter runes."
}

// Error, raised when two or more patterns match with equal match length.
type AmbiguousTokenPattersMatch struct{}

func (self *AmbiguousTokenPattersMatch) Error() string {
	return "Two or more token patterns matched, with same match length."
}
