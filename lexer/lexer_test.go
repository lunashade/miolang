package lexer

import (
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
		if err != nil {
			t.Fatalf("error: %q", err)
		}
		if r == EOF {
			break
		}
		gots = append(gots, r)
	}
	if string(gots) != wants {
		t.Fatalf("want: %q, got: %q", wants, gots)
	}
}
