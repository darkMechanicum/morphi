package lexer

import (
	"github.com/darkMechanicum/morphi/readers"
)

// defaultLexer holds current lexer processing state,
// such as context, current captured string, channels for ready tokens and
// available runes, etc.
type defaultLexer struct {
	// Are there more runes, or last one had been read.
	isEnd bool

	// String that is currently processed and
	// yet hadn't been recognized.
	currentBulk string

	// Available runes for processing
	runeReader readers.RuneReader

	// Lexer config that is currently in use
	cfg *LexerConfig

	// Captured error
	curErr LexerError
}

// Create new default lexer.
func NewDefaultLexer(cfg *LexerConfig, runeReader readers.RuneReader) Lexer {
	return &defaultLexer{false, "", runeReader, cfg, nil}
}

// Read new rune and add it to the current bulk.
func (lexer *defaultLexer) readRune() {
	readRune, err := lexer.runeReader.NextRune()
	switch {
	case err != nil:
		lexer.curErr = err
		return
	case readRune == nil:
		lexer.isEnd = true
	default:
		lexer.currentBulk += string(*readRune)
	}
}

// See interface description.
func (lexer *defaultLexer) NextToken() Token {
	var foundToken Token
	var err error
	switch {
	case lexer.curErr != nil:
		return nil
	case lexer.isEnd && lexer.currentBulk == "":
		return nil
	}

	lexer.readRune()
	if foundToken, err = lexer.greedyMatchPatternToken(); foundToken != nil {
		if err != nil {
			lexer.curErr = err
			return nil
		}
		return foundToken
	}
	return foundToken
}

// See interface description.
func (lexer *defaultLexer) CurrentError() LexerError {
	return lexer.curErr
}

// Try get pattern token.
func (lexer *defaultLexer) greedyMatchPatternToken() (Token, error) {
	for !lexer.isEnd {
		var chosenType TokenType
		var chosenIndex int
		shouldFail := false
		lexer.readRune()
		for _, holder := range lexer.cfg.PatternTokenTypes {
			if interval := holder.pattern.Matches(lexer.currentBulk); interval != nil && interval.start == 0 {
				// If it is not end, then token can match more content.
				if interval.end == len(lexer.currentBulk) && !lexer.isEnd {
					chosenType = nil
					break
				}
				// If not all bulk is matched, then choose largest match.
				if chosenType == nil || chosenIndex < interval.end {
					shouldFail = false
					chosenIndex = interval.end
					chosenType = holder.tType
				} else if chosenIndex == interval.end {
					shouldFail = true
				}
			}
		}
		if shouldFail {
			return nil, &AmbiguousTokenPattersMatch{}
		} else if chosenType != nil {
			tokenContent := lexer.currentBulk[:chosenIndex]
			lexer.currentBulk = lexer.currentBulk[chosenIndex:]
			return &DefaultToken{tokenContent, chosenType}, nil
		}
	}
	return nil, nil
}
