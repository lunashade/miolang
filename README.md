# miolang

yet another whitespace clone

## IMP

- SP -> stack
- TAB SP -> Num Op
- TAB TAB -> Heap
- TAB LF -> IO
- LF -> Flow

## Commands

### stack

- SP <INT>: push INT to stack
- LF SP: copy stack top and push
- LF TAB: swap stack top2
- LF LF: drop stack top
- TAB SP <INT>: copy INT-th value on stack and push it
- TAB LF <INT>: while keeping stack top, drop INT elements on stack

### Num Op

- SP SP: Add s0 + s1
- SP TAB: Sub s0 - s1
- SP LF: Mul s0 * s1
- TAB SP: Div s
- 
