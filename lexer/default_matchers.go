package lexer

import (
	"errors"
	"regexp"
)

type runeSliceDelimiterMatcher struct {
	delimiters map[rune]rune
}

func NewRuneSliceDelimiterMatcher(runes []rune) DelimiterMatcher {
	newMap := make(map[rune]rune)
	for _, rn := range runes {
		newMap[rn] = rn
	}
	return &runeSliceDelimiterMatcher{newMap}
}

func (matcher *runeSliceDelimiterMatcher) doMatch(content string) []int {
	startIndex := -1
	for index, value := range content {
		_, exist := matcher.delimiters[value]
		switch {
		case exist && startIndex == -1:
			startIndex = index
		case !exist && startIndex != -1:
			return []int{startIndex, index}
		}
	}
	return nil
}

func (matcher *runeSliceDelimiterMatcher) MayMatch(content string) []int {
	return matcher.doMatch(content)
}

func (matcher *runeSliceDelimiterMatcher) FullMatch(content string) []int {
	return matcher.doMatch(content)
}

func (matcher *runeSliceDelimiterMatcher) Includes(*DelimiterMatcher) bool {
	return false
}

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

func (pattern *regexpTokenPattern) Matches(content string) int {
	matchIndexes := pattern.regexp.FindStringIndex(content)
	if matchIndexes == nil || matchIndexes[0] != 0 || matchIndexes[0] == matchIndexes[1] {
		return -1
	} else {
		return matchIndexes[1]
	}
}

func (pattern *regexpTokenPattern) Includes(*TokenPattern) bool {
	return false
}

// Create new LexerConfig by baking regexps.
func NewDefaultLexerConfig(
	delimiters []rune,
	fixedTokenTypes map[string]TokenType,
	regexTokenTypes map[string]TokenType,
) (*LexerConfig, error) {
	// bake regexps
	bakedRegexps := make([]TokenPatternHolder, 10)
	for rawRegexp, tType := range regexTokenTypes {
		baked, err := NewRegexpTokenPattern(rawRegexp)
		if err != nil {
			return nil, errors.New("Can't create LexerConfig, because of compiling regexp: " + err.Error())
		}
		holder := TokenPatternHolder{baked, tType}
		bakedRegexps = append(bakedRegexps, holder)
	}
	// bake delimiters
	delimitersMatcher := NewRuneSliceDelimiterMatcher(delimiters)
	// return created context
	return &LexerConfig{
		delimitersMatcher,
		fixedTokenTypes,
		bakedRegexps}, nil
}
