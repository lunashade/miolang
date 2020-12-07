package machine

import (
	"bufio"
	"fmt"
	"io"
	"miolang/commands"
	"os"
	"strconv"
)

type Machine struct {
	s         Stack
	h         Heap
	cmds      []commands.Command
	labels    map[int]int
	callstack Stack
	in        *bufio.Reader
	out       io.Writer
}

func NewMachine(ch chan commands.Command) *Machine {
	cmds := make([]commands.Command, 0)
	labels := make(map[int]int)
	pc := 0
	for cmd := range ch {
		// fmt.Fprintf(os.Stderr, "%d: %q\n", pc, cmd)
		if cmd.Type == commands.IRREGAL {
			panic(fmt.Errorf("IRREGAL: %w", cmd.Err))
		}
		if cmd.Type == commands.FLOW_LABEL {
			labels[cmd.Arg] = pc
		}
		cmds = append(cmds, cmd)
		pc++
	}
	return &Machine{
		s:         make(Stack, 0),
		h:         make(Heap),
		cmds:      cmds,
		labels:    labels,
		callstack: make(Stack, 0),
		in:        bufio.NewReader(os.Stdin),
		out:       os.Stdout,
	}
}

func (m *Machine) Run() {
	for pc := 0; pc < len(m.cmds); {
		cmd := m.cmds[pc]
		// fmt.Fprintf(os.Stderr, "%d: %q\n", pc, cmd)
		switch cmd.Type {
		default:
			panic("IRREGAL")

		// stack
		case commands.STACK_PUSH:
			m.s.Push(cmd.Arg)
		case commands.STACK_DUP:
			m.s.Dup()
		case commands.STACK_SWAP:
			m.s.Swap()
		case commands.STACK_DROP:
			m.s.Drop()
		case commands.STACK_DUP_N:
			m.s.DupN(cmd.Arg)
		case commands.STACK_SLIDE_N:
			m.s.SlideN(cmd.Arg)

		// arith
		case commands.ARITH_ADD:
			m.s.Add()
		case commands.ARITH_SUB:
			m.s.Sub()
		case commands.ARITH_MUL:
			m.s.Mul()
		case commands.ARITH_DIV:
			m.s.Div()
		case commands.ARITH_MOD:
			m.s.Mod()
		// heap
		case commands.HEAP_STORE:
			m.Store_Heap()
		case commands.HEAP_LOAD:
			m.Load_Heap()
		// io
		case commands.IO_PUTCHAR:
			ch := m.s.Pop()
			fmt.Fprintf(m.out, "%c", ch)
		case commands.IO_PUTNUM:
			val := m.s.Pop()
			fmt.Fprintf(m.out, "%d", val)
		case commands.IO_READCHAR:
			r, _, err := m.in.ReadRune()
			if err != nil {
				panic(err)
			}
			addr := m.s.Pop()
			m.h.Store(addr, int(r))
		case commands.IO_READNUM:
			s, err := m.in.ReadString('\n')
			if err != nil {
				panic(err)
			}
			addr := m.s.Pop()
			val, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			m.h.Store(addr, val)
		// flow
		case commands.FLOW_LABEL:
			// do nothing
		case commands.FLOW_GOSUB:
			m.callstack.Push(pc)
			pc = m.labels[cmd.Arg]
		case commands.FLOW_JUMP:
			pc = m.labels[cmd.Arg]
		case commands.FLOW_BEZ:
			val := m.s.Pop()
			if val == 0 {
				pc = m.labels[cmd.Arg]
			}
		case commands.FLOW_BLTZ:
			val := m.s.Pop()
			if val < 0 {
				pc = m.labels[cmd.Arg]
			}
		case commands.FLOW_ENDSUB:
			pc = m.callstack.Pop()
		case commands.FLOW_HALT:
			return
		}
		pc++
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
