package machine

func (s *Stack) Add() {
	a, b := s.Pop(), s.Pop()
	s.Push(a + b)
}
func (s *Stack) Sub() {
	a, b := s.Pop(), s.Pop()
	s.Push(b - a)
}
func (s *Stack) Mul() {
	a, b := s.Pop(), s.Pop()
	s.Push(b * a)
}
func (s *Stack) Div() {
	a, b := s.Pop(), s.Pop()
	s.Push(b / a)
}
func (s *Stack) Mod() {
	a, b := s.Pop(), s.Pop()
	s.Push(b % a)
}
