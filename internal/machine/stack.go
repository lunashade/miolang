package machine

type Stack []int

func (s *Stack) Push(i int) {
	*s = append(*s, i)
}
func (s *Stack) Drop() {
	*s = (*s)[:len(*s)-1]
}
func (s *Stack) Pop() int {
	val := (*s)[len(*s)-1]
	s.Drop()
	return val
}

func (s *Stack) Index(i int) int {
	return (*s)[len(*s)-1-i]
}
func (s *Stack) Dup() { s.DupN(0) }
func (s *Stack) DupN(n int) {
	s.Push(s.Index(n))
}
func (s *Stack) Swap() {
	v1 := s.Pop()
	v2 := s.Pop()
	s.Push(v1)
	s.Push(v2)
}

// [1,2,3,4].SlideN(2) = [1,4]
func (s *Stack) SlideN(n int) {
	val := s.Pop()
	*s = (*s)[:len(*s)-n]
	s.Push(val)
}
