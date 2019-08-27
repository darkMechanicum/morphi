package readers

import (
	"bufio"
	"io"
)

type defaultReader struct {
	inner io.RuneReader
	buff  string
}

func NewDefaultReader(rd io.Reader) RuneReader {
	inner := bufio.NewReader(rd)
	return &defaultReader{inner, ""}
}

func (reader *defaultReader) NextRune() (*rune, error) {
	if reader.buff != "" {
		returnRune := rune(reader.buff[0])
		reader.buff = reader.buff[1:]
		return &returnRune, nil
	} else {
		rn, _, err := reader.inner.ReadRune()
		switch {
		case err == io.EOF:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			return &rn, nil
		}
	}
}

func (reader *defaultReader) PushBack(runes string) {
	reader.buff = reader.buff + runes
}
