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

	// Clean lexer inner state to get fresh token
	// (in case of lexical errors).
	DropBulk()

	// Get next token from reader.
	NextToken() Token

	// Captured error is any.
	CurrentError() LexerError
}

// Abstract Token pattern to determine token type.
// Is represented as interface type since it must
// implement Includes method to exclude conflicting
// patterns.
type TokenPattern interface {
	// Does pattern match passed string.
	// Matches from the beginning.
	// Returns -1 if pattern not found.
	Matches(content string) int

	// Does all matches of current pattern will include
	// matches of passed pattern.
	Includes(*TokenPattern) bool
}

// Struct that holds TokenPattern with its matching TokenType.
type TokenPatternHolder struct {
	pattern TokenPattern
	tType   TokenType
}

// Acts like TokenPattern but matches from the end,
// since we need to detect delimiters as soon as possible.
type DelimiterMatcher interface {
	// Does pattern or its part match passed string.
	// Returns index of delimiter start and end, or nil.
	MayMatch(content string) []int

	// Does pattern fully match passed string.
	// Returns index of delimiter start and end, or nil.
	FullMatch(content string) []int

	// Does all matches of current matcher will include
	// matches of passed matcher.
	Includes(*DelimiterMatcher) bool
}

// Generic lexer config.
type LexerConfig struct {
	// All delimiters, aggregated in single matcher, since
	// they do not produce any token.
	Delimiters DelimiterMatcher

	// Fixed tokens.
	FixedTokenTypes map[string]TokenType

	// Token patterns.
	PatternTokenTypes []TokenPatternHolder
}

// Generic lexer config reader.
type LexerConfigReader interface {

	// Read lexer config from any source.
	ReadConfig(r io.Reader) (*LexerConfig, error)
}
