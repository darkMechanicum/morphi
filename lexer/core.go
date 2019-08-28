package lexer

import (
	"fmt"
	"io"
)

// Generic Token type.
type TokenType interface {
	fmt.Stringer
}

// Abstract lexer's token
type Token interface {
	// Token's string content
	Content() string

	// Token's type
	Type() TokenType
}

// Generic Lexer interface type.
// Can only extract tokens from Reader.
type Lexer interface {
	// Get next token from reader.
	NextToken() Token

	// Captured error is any.
	CurrentError() LexerError
}

// Interval with start (inclusive) and end (exclusive).
type Interval struct {
	start, end int
}

// Abstract Token pattern to determine token type.
// Is represented as interface type since it must
// implement Includes method to exclude conflicting
// patterns.
type TokenPattern interface {
	// Does pattern match passed string.
	// Matches from the beginning.
	// Returns -1 if pattern not found.
	Matches(content string) *Interval
}

// Struct that holds TokenPattern with its matching TokenType.
type TokenPatternHolder struct {
	pattern TokenPattern
	tType   TokenType
}

// Generic lexer config.
type LexerConfig struct {
	// Token patterns.
	PatternTokenTypes []TokenPatternHolder
}

// Generic lexer config reader.
type LexerConfigReader interface {
	// Read lexer config from any source.
	ReadConfig(r io.Reader) (*LexerConfig, error)
}
