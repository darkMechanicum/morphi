package lexer

import (
	"errors"
	"github.com/darkMechanicum/morphi/utils"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Structure of a predefined yml lexer config file
type YmlLexerConfigStructure struct {
	Delimiters map[string]string
	Predefined map[string]string
	Regexp     map[string]string
}

// YmlLexerConfigReader is LexerConfigReader implementation
// for reading yml files. Contains no state.
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
		predefined := make(map[string]TokenType)
		for rawTType, value := range t.Predefined {
			predefined[value] = NewDefaultTokenType(rawTType)
		}
		regexp := make(map[string]TokenType)
		for rawTType, value := range t.Regexp {
			regexp[value] = NewDefaultTokenType(rawTType)
		}

		return NewDefaultLexerConfig(
			delimiters,
			predefined,
			regexp)
	}
}

// Read lexer config from yml file.
func ReadLexerConfigFromYmlFile(path string) (*LexerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	cfgReader := YmlLexerConfigReader{}
	lexerConfig, err := cfgReader.ReadConfig(file)
	if err != nil {
		return nil, err
	}
	return lexerConfig, nil
}
