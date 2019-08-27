package lexer

import "strings"

// DefaultLexer holds current lexer processing state,
// such as context, current captured string, channels for ready tokens and
// available runes, etc.
type DefaultLexer struct {
	// Are there more runes, or last one had been read.
	isEnd bool

	// String that is currently processed and
	// yet hadn't been recognized.
	currentBulk string

	// Available runes for processing
	runeReader RuneReader

	// Lexer config that is currently in use
	cfg LexerConfig

	// Captured error
	curErr LexerError
}

// Determine if delimiter ends at some point at passed string.
func (lexer *DefaultLexer) getDelimiterEndIndex(bulkCandidate string) int {
	// If there is non empty error or end, then return.
	switch {
	case lexer.curErr != nil:
		return -1
	case lexer.isEnd:
		interval := lexer.cfg.Delimiters().FullMatch(bulkCandidate)
		if interval != nil {
			return interval[1]
		} else {
			return 0
		}
	}

	// Perform delimiters matching.
	interval := lexer.cfg.Delimiters().MayMatch(bulkCandidate)
	switch {
	case interval == nil:
		// Can happen, for example, if we read "/" first and threat it like delimiter start
		// but next rune is not "*", so it is not a delimiter anymore.
		return 0
	case interval[0] == 0 && interval[1] == len(bulkCandidate):
		return len(bulkCandidate) // TODO what shall we do when interval[0] is not zero?
	default:
		interval := lexer.cfg.Delimiters().FullMatch(bulkCandidate)
		return interval[1]
	}
}

// Determine if delimiter starts at some point at passed string.
func (lexer *DefaultLexer) getDelimiterStartIndex(bulkCandidate string) int {
	// If there is non empty error or end, then return.
	switch {
	case lexer.curErr != nil:
		return -1
	case lexer.isEnd:
		interval := lexer.cfg.Delimiters().FullMatch(bulkCandidate)
		if interval != nil {
			return interval[0]
		} else {
			return len(bulkCandidate)
		}
	}

	// Perform delimiters matching.
	interval := lexer.cfg.Delimiters().MayMatch(bulkCandidate)
	switch {
	case interval == nil:
		return len(bulkCandidate)
	default:
		interval := lexer.cfg.Delimiters().FullMatch(bulkCandidate)
		return interval[0]
	}
}

// Read next bulk if need. Bulk is non delimiter string from
// rune reader.
func (lexer *DefaultLexer) readBulkIfNeed() {
	// If there is non empty error or bulk, then return.
	switch {
	case lexer.curErr != nil:
		return
	case lexer.currentBulk != "":
		return
	}

	// Skip any heading delimiters if any.
	bulkCandidate := ""
	var endIndex int
	for {
		endIndex = lexer.getDelimiterEndIndex(bulkCandidate)
		if !lexer.isEnd && bulkCandidate != "" && endIndex != len(bulkCandidate) {
			break
		} else {
			readRune, err := lexer.runeReader.NextRune()
			switch {
			case err != nil:
				lexer.curErr = err
				return
			case readRune == nil:
				lexer.isEnd = true
			default:
				bulkCandidate += string(*readRune)
			}
		}
	}

	// Cut delimiter prefix and search for next delimiters.
	bulkCandidate = bulkCandidate[endIndex:]
	var startIndex int
	for {
		startIndex = lexer.getDelimiterStartIndex(bulkCandidate)
		if !lexer.isEnd && bulkCandidate != "" && startIndex != len(bulkCandidate) {
			break
		} else {
			readRune, err := lexer.runeReader.NextRune()
			switch {
			case err != nil:
				lexer.curErr = err
				return
			case readRune == nil:
				lexer.isEnd = true
			default:
				bulkCandidate += string(*readRune)
			}
		}
	}

	// Cut starting delimiter and return it to the reader.
	delimiterStart := bulkCandidate[startIndex:]
	lexer.runeReader.PushBack(delimiterStart)

	// Determine resulting bulk.
	lexer.currentBulk = bulkCandidate[:startIndex-1]
}

// See interface description.
func (lexer *DefaultLexer) DropBulk() {
	switch {
	case lexer.curErr != nil:
		return
	case lexer.isEnd && lexer.currentBulk == "":
		return
	}
	// Force to drop this bulk.
	lexer.currentBulk = ""
	lexer.readBulkIfNeed()
}

// See interface description.
func (lexer *DefaultLexer) NextToken() Token {
	switch {
	case lexer.curErr != nil:
		return nil
	case lexer.isEnd && lexer.currentBulk == "":
		return nil
	}
	lexer.readBulkIfNeed()
	if fixedToken := lexer.greedyMatchFixedToken(); fixedToken != nil {
		return fixedToken
	}
	patternToken, err := lexer.greedyMatchPatternToken()
	if err != nil {
		lexer.curErr = err
		return nil
	}
	return patternToken
}

// See interface description.
func (lexer *DefaultLexer) CurrentError() LexerError {
	return lexer.curErr
}

// Try get fixed token.
func (lexer *DefaultLexer) greedyMatchFixedToken() Token {
	for fixedToken, tType := range lexer.cfg.FixedTokenTypes() {
		if strings.HasPrefix(lexer.currentBulk, fixedToken) {
			lexer.currentBulk = lexer.currentBulk[len(fixedToken):]
			return &DefaultToken{fixedToken, tType}
		}
	}
	return nil
}

// Try get pattern token.
func (lexer *DefaultLexer) greedyMatchPatternToken() (Token, error) {
	var chosenType TokenType
	var chosenIndex int
	for tokenPattern, tType := range lexer.cfg.PatternTokenTypes() {
		if matchIndex := tokenPattern.Matches(lexer.currentBulk); matchIndex > 0 {
			if chosenType == nil || chosenIndex < matchIndex {
				chosenIndex = matchIndex
				chosenType = tType
			} else if chosenIndex == matchIndex {
				return nil, &AmbiguousTokenPattersMatch{}
			}
		}
	}
	if chosenType != nil {
		tokenContent := lexer.currentBulk[:chosenIndex-1]
		lexer.currentBulk = lexer.currentBulk[chosenIndex:]
		return &DefaultToken{tokenContent, chosenType}, nil
	} else {
		return nil, nil
	}
}
