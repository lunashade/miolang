package machine

func Add(s *Stack) {
	a, b := s.Pop(), s.Pop()
	s.Push(a + b)
}
func Sub(s *Stack) {
	a, b := s.Pop(), s.Pop()
	s.Push(b - a)
}
func Mul(s *Stack) {
	a, b := s.Pop(), s.Pop()
	s.Push(b * a)
}
func Div(s *Stack) {
	a, b := s.Pop(), s.Pop()
	s.Push(b / a)
}
func Mod(s *Stack) {
	a, b := s.Pop(), s.Pop()
	s.Push(b % a)
}
