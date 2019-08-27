package lexer

import "errors"

// LexerConfig context is a holder for scope of
// Token describing rules.
type DefaultLexerConfig struct {

	// A set of runes, used as delimiters
	delimiters map[rune]rune

	// Delimiters matcher
	delimitersMatcher DelimiterMatcher

	// A set of fixed tokes
	fixedTokenTypes map[string]*TokenType

	// A set of regex tokes
	regexTokenTypes map[equalAbleRegexp]*TokenType
}

// Create new LexerConfig by baking regexps.
func NewLexerConfig(delimiters []rune, fixedTokenTypes map[string]*TokenType, regexTokenTypes map[string]*TokenType) (*LexerConfig, error) {
	// bake regexps
	bakedRegexps := make(map[equalAbleRegexp]*TokenType)
	for rawRegexp, tType := range regexTokenTypes {
		regexHolder, err := newEqualAbleRegexp(rawRegexp)
		if err != nil {
			return nil, errors.New("Can't create LexerConfig, because of compiling regexp: " + err.Error())
		}
		bakedRegexps[*regexHolder] = tType
	}
	// bake delimiters
	bakedDelimiters := make(map[rune]rune)
	for _, delimiter := range delimiters {
		bakedDelimiters[delimiter] = delimiter
	}
	// return created context
	return &LexerConfig{
		bakedDelimiters,
		fixedTokenTypes,
		bakedRegexps}, nil
}
