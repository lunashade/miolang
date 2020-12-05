package machine

type CmdType int

const (
	IRREGAL CmdType = iota
	// stack SP
	STACK_PUSH    // SP <INT>
	STACK_DUP     // LF SP
	STACK_SWAP    // LF TAB
	STACK_DROP    // LF LF
	STACK_DUP_N   // TAB SP <INT>
	STACK_SLIDE_N // TAB LF <INT>
	// arithmetics TAB SP
	ARITH_ADD // SP SP
	ARITH_SUB // SP TAB
	ARITH_MUL // SP LF
	ARITH_DIV // TAB SP
	ARITH_MOD // TAB TAB
	// heap TAB TAB
	HEAP_STORE // SP
	HEAP_LOAD  // TAB
	// flow: LF
	FLOW_LABEL  // SP SP <UINT>
	FLOW_GOSUB  // SP TAB <UINT>
	FLOW_JUMP   // SP LF <UINT>
	FLOW_BEZ    // TAB SP <UINT>
	FLOW_BLTZ   // TAB TAB <UINT>
	FLOW_ENDSUB // TAB LF
	FLOW_HALT   // LF LF
	// IO: TAB LF
	IO_PUTCHAR  // SP SP
	IO_PUTNUM   // SP TAB
	IO_READCHAR // TAB SP
	IO_READNUM  // TAB TAB
)


type Command struct {
	Type CmdType
	Arg  int
}
