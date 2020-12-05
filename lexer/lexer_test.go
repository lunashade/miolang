package lexer

import (
	"io"
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	inputs := "みおみうくみこっ"
	wants := "みおみみっ"
	l := NewLexer(strings.NewReader(inputs))
	gots := make([]rune, 0, 5)
	for {
		r, err := l.Next()
		if err == io.EOF {
			break
		}
		gots = append(gots, r)
	}
	if string(gots) != wants {
		t.Fatalf("want: %q, got: %q", wants, gots)
	}
}
