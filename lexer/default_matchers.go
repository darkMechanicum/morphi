package lexer

import (
	"errors"
	"regexp"
)

// --- Tokens Token Pattern ---

type regexpTokenPattern struct {
	rawRegexp string
	regexp    *regexp.Regexp
}

func NewRegexpTokenPattern(raw string) (TokenPattern, error) {
	baked, err := regexp.Compile(raw)
	if err != nil {
		return nil, err
	} else {
		return &regexpTokenPattern{raw, baked}, nil
	}
}

func (pattern *regexpTokenPattern) Matches(content string) *Interval {
	matchIndexes := pattern.regexp.FindStringIndex(content)
	if matchIndexes == nil {
		return nil
	} else {
		return &Interval{matchIndexes[0], matchIndexes[1]}
	}
}

func (pattern *regexpTokenPattern) Includes(*TokenPattern) bool {
	return false
}

// Create new LexerConfig by baking regexps.
func NewDefaultLexerConfig(
	regexTokenTypes map[string]TokenType,
) (*LexerConfig, error) {
	// bake regexps
	bakedRegexps := make([]TokenPatternHolder, 0)
	for rawRegexp, tType := range regexTokenTypes {
		baked, err := NewRegexpTokenPattern(rawRegexp)
		if err != nil {
			return nil, errors.New("Can't create LexerConfig, because of compiling regexp: " + err.Error())
		}
		holder := TokenPatternHolder{baked, tType}
		bakedRegexps = append(bakedRegexps, holder)
	}
	// return created context
	return &LexerConfig{bakedRegexps}, nil
}
