package lexer

// Token type is a unique string
// representing this Token semantics.
type DefaultTokenType string

// Stringer implementation for DefaultTokenType.
func (self *DefaultTokenType) String() string {
	return string(*self)
}

// Create and register new Token type.
func NewDefaultTokenType(content string) TokenType {
	result := DefaultTokenType(content)
	return &result
}

// Token is just a string with link to
// LexerConfig.
type DefaultToken struct {

	// Token content.
	myContent string

	// This Token type
	myType TokenType
}

// DefaultToken Token.Content() implementation.
func (self *DefaultToken) Content() string {
	return self.myContent
}

// DefaultToken Token.Type() implementation.
func (self *DefaultToken) Type() TokenType {
	return self.myType
}
