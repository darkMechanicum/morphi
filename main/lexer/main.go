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
// LexerContext.
type Token struct {
	// Token content.
	content string

	// This Token type
	TokenType *TokenType

	// Lexer context link
	context *LexerContext
}

// LexerStateMachine holds current lexer processing state,
// such as context, current rune, channels for ready tokens and
// available runes, etc.
type LexerStateMachine struct {
	// String that is currently processed and
	// yet hadn't been recognized
	current string

	// Available runes for processing
	runes <-chan rune

	// Channel that is used for processed Token passing
	tokens chan<- *Token

	// Lexer context that is currently in use
	ctx *LexerContext
}

func NewLexerStateMachine(runes <-chan rune, ctx *LexerContext) (*LexerStateMachine, <-chan *Token) {
	tokens := make(chan *Token)
	return &LexerStateMachine{"", runes, tokens, ctx}, tokens
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
func (this *equalAbleRegexp) matchString(s string) bool {
	return this.baked.MatchString(s)
}

// Try to extract Token from current string and push it to channel.
// Returns false, if on other tokens are available
func (this *LexerStateMachine) TryGetToken() (bool, error) {
	var ok bool
	var nextRune rune
	for nextRune, ok = <-this.runes; ok; nextRune, ok = <-this.runes {
		log.Printf("Read next rune <%s>\n", string(nextRune))
		if _, exists := this.ctx.delimiters[nextRune]; exists {
			log.Printf("It is delimiter <%s>\n", string(nextRune))
			break
		} else {
			this.current += string(nextRune)
		}
	}
	if !ok {
		log.Printf("Channel closed. Last string is <%s>", this.current)
		if len(this.current) > 0 {
			return this.determineAndPush()
		}
		close(this.tokens)
		return false, nil
	}
	return this.determineAndPush()
}

func (this *LexerStateMachine) determineAndPush() (bool, error) {
	if chosenType, exist := this.ctx.fixedTokenTypes[this.current]; exist {
		log.Printf("Chosen fixed token type <%s> for <%s>", string(*chosenType), this.current)
		this.pushToken(chosenType)
		return true, nil
	}
	var chosenType *TokenType
	for curRegexp, tType := range this.ctx.regexTokenTypes {
		if curRegexp.matchString(this.current) {
			if chosenType == nil {
				chosenType = tType
			} else {
				return true, errors.New("Two regexp matched the Token.")
			}
		}
	}
	if chosenType == nil {
		log.Printf("Chosen default token type <%s> for <%s>", string(*this.ctx.defaultTokenType), this.current)
		this.pushToken(this.ctx.defaultTokenType)
		return true, nil
	} else {
		log.Printf("Chosen regexp token type <%s> for <%s>", string(*chosenType), this.current)
		this.pushToken(chosenType)
		return true, nil
	}
}

// Create Token from current string with passed type and push it to channel.
// Also, clears current string.
func (this *LexerStateMachine) pushToken(tType *TokenType) {
	newToken := &Token{this.current, tType, this.ctx}
	log.Print("Pushing new token", newToken)
	this.tokens <- newToken
	this.current = ""
}

// LexerContext context is a holder for scope of
// Token describing rules.
type LexerContext struct {
	// A set of runes, used as delimiters
	delimiters map[rune]rune

	// A set of fixed tokes
	fixedTokenTypes map[string]*TokenType

	// A set of regex tokes
	regexTokenTypes map[equalAbleRegexp]*TokenType

	// Default Token
	defaultTokenType *TokenType
}

// Create new LexerContext by baking regexps.
func NewLexerContext(delimiters []rune, fixedTokenTypes map[string]*TokenType, regexTokenTypes map[string]*TokenType, defaultTokenType *TokenType) (*LexerContext, error) {
	// bake regexps
	bakedRegexps := make(map[equalAbleRegexp]*TokenType)
	for rawRegexp, tType := range regexTokenTypes {
		bakedRegexp, err := regexp.Compile(rawRegexp)
		if err != nil {
			return nil, errors.New("Can't create LexerContext, because of compiling regexp: " + err.Error())
		}
		regexHolder := equalAbleRegexp{rawRegexp, bakedRegexp}
		bakedRegexps[regexHolder] = tType
	}
	// bake delimiters
	bakedDelimiters := make(map[rune]rune)
	for _, delimiter := range delimiters {
		bakedDelimiters[delimiter] = delimiter
	}
	// return created context
	return &LexerContext{
		bakedDelimiters,
		fixedTokenTypes,
		bakedRegexps,
		defaultTokenType}, nil
}
