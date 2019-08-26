package main

import (
	"./lexer"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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

	// Initialize Context
	cfgReader := lexer.YmlLexerConfigReader{}
	file, err := os.Open("C:/workspace/projects/go/src/github.com/darkMechanicum/morphi/resources/simple_lex.yml")
	if err != nil {
		panic(err)
	}
	lexerConfig, err := cfgReader.ReadConfig(file)
	if err != nil {
		panic(err)
	}

	// Initialize the reader
	input := strings.NewReader("one + three = four <?>")
	runes := startReading(input)

	// Initialize StateMachine
	stateMachine, tokens := lexer.NewLexerStateMachine(runes, lexerConfig)
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
			fmt.Printf("%s (%s)\n", string(*token.TokenType), *token.Content)
		}
	}()
}
