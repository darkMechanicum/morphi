package lexer

import (
	"errors"
	"github.com/darkMechanicum/morphi/utils"
	"gopkg.in/yaml.v2"
	"io"
)

type LexerConfigReader interface {
	ReadConfig(r io.Reader) (*LexerConfig, error)
}

// Structure of a predefined yml lexer config file
type YmlLexerConfigStructure struct {
	Delimiters map[string]string
	Predefined map[string]string
	Regexp     map[string]string
}

type YmlLexerConfigReader struct {
}

func (self *YmlLexerConfigReader) ReadConfig(r io.Reader) (*LexerConfig, error) {
	t := &YmlLexerConfigStructure{}
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(t)
	if err != nil {
		return nil, err
	} else {
		delimiters := make([]rune, 0, len(t.Delimiters))
		for _, delimiterStr := range t.Delimiters {
			delimiter, err := utils.GetOnlyRune(&delimiterStr)
			if err != nil {
				return nil, errors.New("Can't read lexer config because delimiter string is empty " +
					"or contains more than one character.")
			}
			delimiters = append(delimiters, delimiter)
		}
		predefined := make(map[string]*TokenType)
		for rawTType, value := range t.Predefined {
			predefined[value] = NewTokenType(rawTType)
		}
		regexp := make(map[string]*TokenType)
		for rawTType, value := range t.Regexp {
			regexp[value] = NewTokenType(rawTType)
		}

		return NewLexerConfig(
			delimiters,
			predefined,
			regexp)
	}
}
