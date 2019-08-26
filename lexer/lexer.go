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
func (self *LexerStateMachine) TryGetToken() (bool, error) {
	for nextRune, canNext := <-self.runes; canNext; nextRune, canNext = <-self.runes {
		log.Printf("Read next rune <%s>\n", string(nextRune))
		self.addToCurrent(nextRune)
		if _, exists := self.ctx.delimiters[nextRune]; exists {
			log.Printf("It is delimiter <%s>\n", string(nextRune))
			return self.determineAndPush(true)
		} else {
			pushed, err := self.determineAndPush(false)
			if err != nil {
				return false, err
			} else if pushed {
				return true, nil
			}
		}
	}

	log.Printf("Channel closed. Last string is <%s>", *self.current)
	_, err := self.determineAndPush(true)
	close(self.tokens)
	return false, err
}

// Add next rune to the current string, ignoring delimiters.
func (self *LexerStateMachine) addToCurrent(nextRune rune) {
	if _, exists := self.ctx.delimiters[nextRune]; !exists {
		*self.current += string(nextRune)
	}
}

// Try to determine token type for current string
// and push it in token channel.
func (self *LexerStateMachine) determineAndPush(isEnd bool) (bool, error) {
	if len(*self.current) == 0 {
		return true, nil
	}
	if chosenType, exist := self.ctx.fixedTokenTypes[*self.current]; exist {
		log.Printf("Chosen fixed token type <%s> for <%s>", string(*chosenType), *self.current)
		self.pushToken(chosenType, nil)
		return true, nil
	}
	if isEnd {
		return self.determineRegexpAndPush()
	} else {
		return true, nil
	}
}

// Try to determine regexp token type for current string
// and push it in token channel.
func (self *LexerStateMachine) determineRegexpAndPush() (bool, error) {
	var chosenType *TokenType
	var matchedIndex, chosenIndex []int
	for {
		for curRegexp, tType := range self.ctx.regexTokenTypes {
			log.Printf("Trying regexp <%s> for <%s>", curRegexp.rawRegexp, *self.current)
			matchedIndex = curRegexp.baked.FindStringIndex(*self.current)
			if matchedIndex != nil && matchedIndex[0] != matchedIndex[1] {
				log.Printf("Regexp <%s> matched for <%s>", curRegexp.rawRegexp, *self.current)
				log.Printf("Matched index is %d-%d", matchedIndex[0], matchedIndex[1])
				if chosenType == nil {
					chosenType = tType
					chosenIndex = matchedIndex
				} else {
					log.Printf("Two regexp matched the Token (last is <%s>).", curRegexp.rawRegexp)
					return false, errors.New("Two regexp matched the Token.")
				}
			}
		}
		if chosenType == nil {
			log.Printf("No token found for <%s>", *self.current)
			return false, nil
		} else {
			if chosenIndex != nil {
				log.Printf("Chosen regexp token type <%s> for <%s>", string(*chosenType), (*self.current))
				if newStart := chosenIndex[1]; newStart < len(*self.current) {
					self.pushToken(chosenType, &newStart)
					return self.determineAndPush(true)
				} else {
					self.pushToken(chosenType, nil)
					log.Printf("Success return from regexp matching")
					return true, nil
				}
			} else {
				log.Printf("Unsuccess return from regexp matching")
				return false, nil
			}
		}
	}
}

// Create Token from current string with passed type and push it to channel.
// Also, clears current string.
func (self *LexerStateMachine) pushToken(tType *TokenType, newStart *int) {
	var newToken *Token
	if newStart == nil {
		log.Printf("Lol")
		newToken = &Token{self.current, tType, self.ctx}
		newCurrent := ""
		self.current = &newCurrent
	} else {
		tokenContent := (*self.current)[:*newStart]
		newToken = &Token{&tokenContent, tType, self.ctx}
		newCurrent := (*self.current)[*newStart:]
		self.current = &newCurrent
	}
	self.tokens <- newToken
	log.Printf("Pushed new token with content <%s> and type <%s>", *newToken.Content, *newToken.TokenType)
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
