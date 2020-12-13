package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	SP  = 'み'
	TAB = 'お'
	LF  = 'っ'
	EOF = -1
)

type Lexer struct {
	*bufio.Reader
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{bufio.NewReader(r)}
}

func (l *Lexer) Next() (rune, error) {
	for {
		r, _, err := l.ReadRune()
		if err == io.EOF {
			return EOF, nil
		}
		if err != nil {
			return EOF, err
		}
		if r == utf8.RuneError {
			return EOF, fmt.Errorf("RuneError")
		}
		switch r {
		case SP, TAB, LF:
			return r, nil
		}
	}
}
