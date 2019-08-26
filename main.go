package main

import (
	"fmt"
	"github.com/darkMechanicum/morphi/lexer"
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
	runes := utils.FromReader(sample)
	if err != nil {
		panic(err)
	}

	// Initialize StateMachine
	stateMachine, tokens := lexer.NewLexerStateMachine(runes, lexerConfig)
	startWriting(tokens)
	for cnt, err := stateMachine.TryGetToken(); cnt; cnt, err = stateMachine.TryGetToken() {
		if err != nil {
			panic(err)
		}
	}
}

func startWriting(tokens <-chan *lexer.Token) {
	go func() {
		for token := range tokens {
			fmt.Printf("%s (%s)\n", string(*token.TokenType), *token.Content)
		}
	}()
}
