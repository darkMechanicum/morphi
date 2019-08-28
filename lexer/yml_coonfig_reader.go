package lexer

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Structure of a predefined yml lexer config file
type YmlLexerConfigStructure struct {
	Tokens map[string]string
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
		regexp := make(map[string]TokenType)
		for rawTType, value := range t.Tokens {
			regexp[value] = NewDefaultTokenType(rawTType)
		}

		return NewDefaultLexerConfig(regexp)
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
