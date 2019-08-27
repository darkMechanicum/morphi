package main

import (
	"fmt"
	"github.com/darkMechanicum/morphi/lexer"
	"github.com/darkMechanicum/morphi/readers"
	"github.com/darkMechanicum/morphi/utils"
	"os"
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
	lexer := lexer.NewDefaultLexer(lexerConfig, runeReader)
	for token := lexer.NextToken(); token != nil; token = lexer.NextToken() {
		fmt.Printf("%s (%s)\n", token.Type().String(), token.Content())
	}
	if lexer.CurrentError() != nil {
		panic(lexer.CurrentError())
	}
}
