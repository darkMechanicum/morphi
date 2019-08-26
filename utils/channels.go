package utils

import (
	"io"
	"text/scanner"
)

// Create a channel and launch goroutine that
// reads runes from passed reader and pushes them into the channel.
func FromReader(reader io.Reader) <-chan rune {
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
