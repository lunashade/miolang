package parser

import (
	"fmt"
	"io"
	"miolang/commands"
	"miolang/lexer"
)

type Parser struct {
	l        *lexer.Lexer
	Commands chan commands.Command
	err      error
}

var (
	UnexpectedEOF  = fmt.Errorf("unexpected EOF")
	UnknownCommand = fmt.Errorf("unknown commands")
)

func (p *Parser) Next() rune {
	r, err := p.l.Next()
	if err != nil {
		panic(err)
	}
	return r
}

func Parse(r io.Reader) (*Parser, chan commands.Command) {
	p := &Parser{
		l:        lexer.NewLexer(r),
		Commands: make(chan commands.Command),
	}
	go p.run()
	return p, p.Commands
}
func (p *Parser) run() {
	for state := parseIMP; state != nil; {
		state = state(p)
	}
	close(p.Commands)
}
func (p *Parser) emit(ty commands.CmdType) {
	p.Commands <- commands.Command{Type: ty}
}
func (p *Parser) emitError(err error) {
	p.Commands <- commands.Command{Type: commands.IRREGAL, Err: err}
}
func (p *Parser) emitVal(ty commands.CmdType) {
	val, err := parseInt(p)
	if err != nil {
		p.emitError(err)
	} else {
		p.Commands <- commands.Command{Type: ty, Arg: val}
	}
}
func (p *Parser) emitLabel(ty commands.CmdType) {
	val, err := parseUInt(p)
	if err != nil {
		p.emitError(err)
	} else {
		p.Commands <- commands.Command{Type: ty, Arg: val}
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
			p.emitError(UnexpectedEOF)
			return nil
		case lexer.SP:
			return parseArith
		case lexer.TAB:
			return parseHeap
		case lexer.LF:
			return parseIO
		}
	default:
		p.emitError(UnexpectedEOF)
		return nil
	}
}

func parseStack(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emitError(UnexpectedEOF)
		return nil
	case lexer.SP:
		p.emitVal(commands.STACK_PUSH)
		return parseIMP
	case lexer.LF:
		r2 := p.Next()
		switch r2 {
		case lexer.SP:
			p.emit(commands.STACK_DUP)
		case lexer.TAB:
			p.emit(commands.STACK_SWAP)
		case lexer.LF:
			p.emit(commands.STACK_DROP)
		default:
			p.emitError(fmt.Errorf("parseStack: %w", UnexpectedEOF))
			return nil
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		case lexer.SP:
			p.emitVal(commands.STACK_DUP_N)
		case lexer.LF:
			p.emitVal(commands.STACK_SLIDE_N)
		case lexer.TAB:
			p.emitError(fmt.Errorf("parseStack: %w", UnknownCommand))
			return nil
		default:
			p.emitError(fmt.Errorf("parseStack: %w", UnexpectedEOF))
			return nil
		}
		return parseIMP
	}
}
func parseArith(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emitError(fmt.Errorf("parseArith: %w", UnknownCommand))
		return nil
	case lexer.SP:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseArith: %w", UnexpectedEOF))
			return nil
		case lexer.SP: // SP SP
			p.emit(commands.ARITH_ADD)
		case lexer.TAB: // SP TAB
			p.emit(commands.ARITH_SUB)
		case lexer.LF: // SP TAB
			p.emit(commands.ARITH_MUL)
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseArith: %w", UnknownCommand))
			return nil
		case lexer.SP: // TAB SP
			p.emit(commands.ARITH_DIV)
		case lexer.TAB: // TAB TAB
			p.emit(commands.ARITH_MOD)
		}
		return parseIMP
	}
}
func parseFlow(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emitError(fmt.Errorf("parseFlow: %w", UnknownCommand))
		return nil
	case lexer.LF:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseFlow: %w", UnknownCommand))
			return nil
		case lexer.LF:
			p.emit(commands.FLOW_HALT)
		}
		return parseIMP
	case lexer.SP:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseFlow: %w", UnknownCommand))
			return nil
		case lexer.SP:
			p.emitLabel(commands.FLOW_LABEL)
		case lexer.TAB:
			p.emitLabel(commands.FLOW_GOSUB)
		case lexer.LF:
			p.emitLabel(commands.FLOW_JUMP)
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseFlow: %w", UnknownCommand))
			return nil
		case lexer.SP:
			p.emitLabel(commands.FLOW_BEZ)
		case lexer.TAB:
			p.emitLabel(commands.FLOW_BLTZ)
		case lexer.LF:
			p.emit(commands.FLOW_ENDSUB)
		}
		return parseIMP
	}
}
func parseHeap(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emitError(fmt.Errorf("parseHeap: %w", UnknownCommand))
		return nil
	case lexer.SP:
		p.emit(commands.HEAP_STORE)
	case lexer.TAB:
		p.emit(commands.HEAP_LOAD)
	}
	return parseIMP
}
func parseIO(p *Parser) stateFn {
	r := p.Next()
	switch r {
	default:
		p.emitError(fmt.Errorf("parseIO: %w", UnknownCommand))
		return nil
	case lexer.SP:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseIO: %w", UnknownCommand))
			return nil
		case lexer.SP:
			p.emit(commands.IO_PUTCHAR)
		case lexer.TAB:
			p.emit(commands.IO_PUTNUM)
		}
		return parseIMP
	case lexer.TAB:
		r2 := p.Next()
		switch r2 {
		default:
			p.emitError(fmt.Errorf("parseIO: %w", UnknownCommand))
			return nil
		case lexer.SP:
			p.emit(commands.IO_READCHAR)
		case lexer.TAB:
			p.emit(commands.IO_READNUM)
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
// first item parsed as sign, SP for +, TAB for -
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
