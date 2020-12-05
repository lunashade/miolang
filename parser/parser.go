package parser

import (
	"io"
	"miolang/lexer"
	"miolang/machine"
)

type Parser struct {
	l   *lexer.Lexer
	IRs chan machine.IR
	err error
}

func (p *Parser) Next() rune {
	r, err := p.l.Next()
	if err != nil {
		panic(err)
	}
	return r
}

func Parse(r io.Reader) (*Parser, chan machine.IR) {
	p := &Parser{
		l:   lexer.NewLexer(r),
		IRs: make(chan machine.IR),
	}
	go p.run()
	return p, p.IRs
}
func (p *Parser) run() {
	for state := parseIMP; state != nil; {
		state = state(p)
	}
	close(p.IRs)
}
func (p *Parser) emit(ty machine.IRType) {
	p.IRs <- machine.IR{Type: ty}
}
func (p *Parser) emitVal(ty machine.IRType) {
	val, err := parseInt(p)
	if err != nil {
		p.emit(machine.IRREGAL)
	} else {
		p.IRs <- machine.IR{Type: ty, Arg: val}
	}
}
func (p *Parser) emitLabel(ty machine.IRType) {
	val, err := parseUInt(p)
	if err != nil {
		p.emit(machine.IRREGAL)
	} else {
		p.IRs <- machine.IR{Type: ty, Arg: val}
	}
}

type stateFn func(*Parser) stateFn

func parseIMP(p *Parser) stateFn {
	r := p.Next()
	switch r {
	case lexer.EOF:
		return nil
	case lexer.SP:
		return parseStack
	case lexer.LF:
		return parseFlow
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP:
			return parseArith
		case lexer.TAB:
			return parseHeap
		case lexer.LF:
			return parseIO
		}
	default:
		p.emit(machine.IRREGAL)
		return nil
	}
}

func parseStack(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emit(machine.IRREGAL)
		return nil
	case lexer.SP:
		p.emitVal(machine.STACK_PUSH)
		return parseIMP
	case lexer.LF:
		r2 := p.Next()
		switch r2 {
		case lexer.SP:
			p.emit(machine.STACK_DUP)
		case lexer.TAB:
			p.emit(machine.STACK_SWAP)
		case lexer.LF:
			p.emit(machine.STACK_DROP)
		default:
			p.emit(machine.IRREGAL)
			return nil
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		case lexer.SP:
			p.emitVal(machine.STACK_DUP_N)
		case lexer.LF:
			p.emitVal(machine.STACK_SLIDE_N)
		default:
			p.emit(machine.IRREGAL)
			return nil
		}
		return parseIMP
	}
}
func parseArith(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emit(machine.IRREGAL)
		return nil
	case lexer.SP:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP: // SP SP
			p.emit(machine.ARITH_ADD)
		case lexer.TAB: // SP TAB
			p.emit(machine.ARITH_SUB)
		case lexer.LF: // SP TAB
			p.emit(machine.ARITH_MUL)
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP: // TAB SP
			p.emit(machine.ARITH_DIV)
		case lexer.TAB: // TAB TAB
			p.emit(machine.ARITH_MOD)
		}
		return parseIMP
	}
}
func parseFlow(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emit(machine.IRREGAL)
		return nil
	case lexer.LF:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.LF:
			p.emit(machine.FLOW_HALT)
			return nil
		}
	case lexer.SP:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP:
			p.emitLabel(machine.FLOW_LABEL)
		case lexer.TAB:
			p.emitLabel(machine.FLOW_GOSUB)
		case lexer.LF:
			p.emitLabel(machine.FLOW_JUMP)
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP:
			p.emitLabel(machine.FLOW_BEZ)
		case lexer.TAB:
			p.emitLabel(machine.FLOW_BLTZ)
		case lexer.LF:
			p.emit(machine.FLOW_ENDSUB)
		}
		return parseIMP
	}
}
func parseHeap(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emit(machine.IRREGAL)
		return nil
	case lexer.SP:
		p.emit(machine.HEAP_STORE)
	case lexer.TAB:
		p.emit(machine.HEAP_LOAD)
	}
	return parseIMP
}
func parseIO(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emit(machine.IRREGAL)
		return nil
	case lexer.SP:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP:
			p.emitLabel(machine.IO_PUTCHAR)
		case lexer.TAB:
			p.emitLabel(machine.IO_PUTNUM)
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emit(machine.IRREGAL)
			return nil
		case lexer.SP:
			p.emitLabel(machine.IO_READCHAR)
		case lexer.TAB:
			p.emitLabel(machine.IO_READNUM)
		}
		return parseIMP
	}
}

// parseUInt parses unsigned Integer.
// SP is 0, TAB is 1, continues until LF is given
// Example:
//    TAB TAB SP LF = 0b110 = 6
func parseUInt(p *Parser) (int, error) {
	val := 0
	for {
		r := p.Next()
		switch r {
		case lexer.LF:
			return val, nil
		case lexer.SP:
			val *= 2
		case lexer.TAB:
			val *= 2
			val += 1
		case lexer.EOF:
			return 0, io.EOF
		}
	}

}

// parseInt parses signed Integer.
// fmachinest item parsed as sign, SP for +, TAB for -
// if LF is given for sign, return 0 without error
func parseInt(p *Parser) (int, error) {
	sign := 1
	r := p.Next()
	switch r {
	case lexer.LF:
		return 0, nil
	case lexer.EOF:
		return 0, io.EOF
	case lexer.SP:
		sign = 1
	case lexer.TAB:
		sign = -1
	}
	val, err := parseUInt(p)
	return val * sign, err
}
