package main

import (
	"fmt"
	"github.com/darkMechanicum/morphi/lexer"
	"github.com/darkMechanicum/morphi/readers"
	"github.com/darkMechanicum/morphi/utils"
	"os"
	"strings"
)

func main() {
	utils.InitLogging()

	// Initialize Context
	lexerConfig, err := lexer.ReadLexerConfigFromYmlFile("C:/workspace/projects/go/src/github.com/darkMechanicum/morphi/resources/simple_lex.yml")
	if err != nil {
		panic(err)
	}

	// Initialize the reader
	sample, err := os.Open("C:/workspace/projects/go/src/github.com/darkMechanicum/morphi/resources/sample.txt")
	runeReader := readers.NewDefaultReader(sample)
	if err != nil {
		panic(err)
	}

	// Initialize StateMachine
	lex := lexer.NewDefaultLexer(lexerConfig, runeReader)
	for token := lex.NextToken(); token != nil; token = lex.NextToken() {
		escapedContent := strings.ReplaceAll(token.Content(), "\r", "")
		escapedContent = strings.ReplaceAll(escapedContent, "\n", "")
		fmt.Printf("%s (%s)\n", token.Type().String(), escapedContent)
	}
	if lex.CurrentError() != nil {
		panic(lex.CurrentError())
	}
}
