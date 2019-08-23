package main

import (
	"../lexer"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"text/scanner"
)

func main() {
	logEnabled := flag.Bool("l", false, "enables logging")
	flag.Parse()

	if !*logEnabled {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	literal := lexer.NewTokenType("literal")
	plus := lexer.NewTokenType("plus")
	equal := lexer.NewTokenType("equal")
	tag := lexer.NewTokenType("tag")

	// Initializae Context
	lexerContext, err := lexer.NewLexerContext(
		[]rune{' '},
		map[string]*lexer.TokenType{"+": plus, "=": equal},
		map[string]*lexer.TokenType{"<.>": tag},
		literal)

	if err != nil {
		panic(err)
	}

	// Initialize the reader
	input := strings.NewReader("one + three = four <?>")
	runes := startReading(input)

	// Initialize StateMachine
	stateMachine, tokens := lexer.NewLexerStateMachine(runes, lexerContext)
	startWriting(tokens)
	for cnt, err := stateMachine.TryGetToken(); cnt; cnt, err = stateMachine.TryGetToken() {
		if err != nil {
			panic(err)
		}
	}
}

func startReading(reader io.Reader) <-chan rune {
	readerChan := make(chan rune) // explicitly set buffer size to 1
	go func() {
		var s scanner.Scanner
		s.Init(reader)
		for current := s.Next(); current != scanner.EOF; current = s.Next() {
			readerChan <- current
		}
		close(readerChan)
	}()
	return readerChan
}

func startWriting(tokens <-chan *lexer.Token) {
	go func() {
		for token := range tokens {
			fmt.Printf("%s\n", string(*token.TokenType))
		}
	}()
}
