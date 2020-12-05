package machine

import (
	"fmt"
	"miolang/commands"
)

type Machine struct {
	s    Stack
	h    Heap
	cmds chan commands.Command
}

func NewMachine(ch chan commands.Command) *Machine {
	return &Machine{
		s:    make(Stack, 0),
		h:    make(Heap),
		cmds: ch,
	}
}

func (m *Machine) Run() {
	for cmd := range m.cmds {
		if cmd.Type == commands.IRREGAL {
			panic("IRREGAL")
		}
		if cmd.Type == commands.FLOW_HALT {
			return
		}
		fmt.Printf("%q", cmd)
	}
}

func (m *Machine) Store_Heap() {
	value, addr := m.s.Pop(), m.s.Pop()
	m.h.Store(addr, value)
}
func (m *Machine) Load_Heap() {
	addr := m.s.Pop()
	m.s.Push(m.h.Load(addr))
}
