package machine

import "testing"

func TestStackSlideN(t *testing.T) {
	s := Stack{1, 2, 3, 4}
	s.SlideN(2)
	want := []int{1, 4}
	for i, got := range s {
		if want[i] != got {
			t.Fatalf("stack slideN fail at %d: want: %q, got: %q",
				i, want, s)
		}
	}
}

func TestStackDupN(t *testing.T) {
	s := Stack{}
	s.Push(1)
	s.Push(2)
	s.DupN(1)
	want := []int{1, 2, 1}
	for i, got := range s {
		if want[i] != got {
			t.Fatalf("stack slideN fail at %d: want: %q, got: %q",
				i, want, s)
		}
	}
}
