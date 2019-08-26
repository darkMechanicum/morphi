package lexer

import (
	"errors"
	"log"
	"regexp"
)

// Token type is a unique string
// representing this Token semantics
type TokenType string

// Create and register new Token type
func NewTokenType(content string) *TokenType {
	result := TokenType(content)
	return &result
}

// Token is just a string with link to
// LexerConfig.
type Token struct {
	// Token content.
	Content *string

	// This Token type
	TokenType *TokenType

	// Lexer context link
	context *LexerConfig
}

// LexerStateMachine holds current lexer processing state,
// such as context, current rune, channels for ready tokens and
// available runes, etc.
type LexerStateMachine struct {
	// String that is currently processed and
	// yet hadn't been recognized
	current *string

	// Available runes for processing
	runes <-chan rune

	// Channel that is used for processed Token passing
	tokens chan<- *Token

	// Lexer context that is currently in use
	ctx *LexerConfig
}

func NewLexerStateMachine(runes <-chan rune, ctx *LexerConfig) (*LexerStateMachine, <-chan *Token) {
	tokens := make(chan *Token)
	newCurrent := ""
	return &LexerStateMachine{&newCurrent, runes, tokens, ctx}, tokens
}

// Regular expression that has backing string.
type equalAbleRegexp struct {
	// String regexp representation.
	rawRegexp string

	// Baked regexp
	baked *regexp.Regexp
}

// Try to bake regexp
func newEqualAbleRegexp(rawRegexp string) (*equalAbleRegexp, error) {
	baked, bakingError := regexp.Compile(rawRegexp)
	if bakingError != nil {
		return nil, bakingError
	}
	return &equalAbleRegexp{rawRegexp, baked}, nil
}

// Try to match regexp
func (this *equalAbleRegexp) matchString(s *string) bool {
	return this.baked.MatchString(*s)
}

// Try to extract Token from current string and push it to channel.
// Returns false, if on other tokens are available
func (this *LexerStateMachine) TryGetToken() (bool, error) {
	for nextRune, canNext := <-this.runes; canNext; nextRune, canNext = <-this.runes {
		log.Printf("Read next rune <%s>\n", string(nextRune))
		this.addToCurrent(nextRune)
		if _, exists := this.ctx.delimiters[nextRune]; exists {
			log.Printf("It is delimiter <%s>\n", string(nextRune))
			return this.determineAndPush(true)
		} else {
			pushed, err := this.determineAndPush(false)
			if err != nil {
				return false, err
			} else if pushed {
				return true, nil
			}
		}
	}

	log.Printf("Channel closed. Last string is <%s>", *this.current)
	_, err := this.determineAndPush(true)
	close(this.tokens)
	return false, err
}

// Add next rune to the current string, ignoring delimiters.
func (this *LexerStateMachine) addToCurrent(nextRune rune) {
	if _, exists := this.ctx.delimiters[nextRune]; !exists {
		*this.current += string(nextRune)
	}
}

// Try to determine token type for current string
// and push it in token channel.
func (this *LexerStateMachine) determineAndPush(isEnd bool) (bool, error) {
	if len(*this.current) == 0 {
		return true, nil
	}
	if chosenType, exist := this.ctx.fixedTokenTypes[*this.current]; exist {
		log.Printf("Chosen fixed token type <%s> for <%s>", string(*chosenType), *this.current)
		this.pushToken(chosenType)
		return true, nil
	}
	var chosenType *TokenType
	if isEnd {
		for curRegexp, tType := range this.ctx.regexTokenTypes {
			if curRegexp.matchString(this.current) {
				if chosenType == nil {
					chosenType = tType
				} else {
					return true, errors.New("Two regexp matched the Token.")
				}
			}
		}
	}
	if chosenType == nil {
		log.Printf("No token found for <%s>", *this.current)
		return false, nil
	} else {
		log.Printf("Chosen regexp token type <%s> for <%s>", string(*chosenType), *this.current)
		this.pushToken(chosenType)
		return true, nil
	}
}

// Create Token from current string with passed type and push it to channel.
// Also, clears current string.
func (this *LexerStateMachine) pushToken(tType *TokenType) {
	newToken := &Token{this.current, tType, this.ctx}
	this.tokens <- newToken
	log.Printf("Pushed new token with content <%s> and type <%s>", *newToken.Content, *newToken.TokenType)
	newCurrent := ""
	this.current = &newCurrent
}

// LexerConfig context is a holder for scope of
// Token describing rules.
type LexerConfig struct {
	// A set of runes, used as delimiters
	delimiters map[rune]rune

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
