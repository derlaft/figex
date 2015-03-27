package asm

import (
    //. "github.com/derlaft/figex/mio"
    //"strconv"
)

type Args struct {

    //Op Instruction
    Op string
    A [2]Arg
    Used byte
}


type State struct {
    Reg [16]byte
    Mem [256]byte
    Ret [32]int
    IP, pt int
}

type pargs struct {
    op Instruction
    A [2]byte
    jump int
    R *byte
}

type Arg struct {
    Type byte
    Val int
}

const (
    ARG_REG = 0
    ARG_CONST = 1
    ARG_LABEL = 2
)

// links to instruction subroutines
var Ops = map[string]Instruction{
    "ADD": Add,
    "SUB": Sub,
    "INC": Inc,
    "DEC": Dec,
    "MUL": Mul,
    "DIV": Div,
    "AND": And,
    "OR":  Or,
    "XOR": Xor,
    "NOT": Not,
    "ROL": Rol,
    "ROR": Ror,
    "RCL": Rcl,
    "RCR": Rcr,
    "MOV": Mov,
    "LD":  Ld,
    "ST":  St,
    "PUT": Push,
    "POP": Pop,
    "JMP": Jmp,
    "JZ":  Jz,
    "JNZ": Jnz,
    "JO":  Jo,
    "JNO": Jno,
    "JF":  Jf,
    "JNF": Jnf,
    "FNC": Call,
    "FLT": Flt,
}

// instruction that need Args.R pointer to return value
var Returning = map[string]bool {
    "ADD": true,
    "SUB": true,
    "INC": true,
    "DEC": true,
    "MUL": true,
    "DIV": true,
    "AND": true,
    "OR":  true,
    "XOR": true,
    "NOT": true,
    "ROL": true,
    "ROR": true,
    "RCL": true,
    "RCR": true,
    "MOV": true,
    "LD":  true,
    "POP": true,
}

const (
  RA = 10
  RB = 11
  RC = 12
  RD = 13
  RE = 14 //stack pointer
  RF = 15 //flag register

  F_ZERO = 0
  F_OVER = 1
  F_FAULT = 2
  F_INT = 3
)

type Instruction func(*State,pargs)

// @TODO: check for interrupts
// and implement them ofc :)

func (state *State) Tick(a Args) error {

    state.Reg[RF] = 0
    pargs := a.getPargs(state)
    pargs.op(state, pargs)
    state.IP += 1

    return nil
}

func (state *State) GetIP() int {
    return state.IP
}

func (a *Args) getPargs(s *State) pargs {
    p := pargs{}

    p.op = Ops[a.Op]

    for i, arg := range a.A {
        switch arg.Type {
            case ARG_REG:
                reg := arg.Val & 0xF
                p.A[i] = s.Reg[reg]
                if i == 0 {
                    p.R = &s.Reg[reg]
                }
            case ARG_CONST:
                p.A[i] = byte(arg.Val)
            case ARG_LABEL:
                p.jump = arg.Val
        }
    }


    if p.R == nil && Returning[a.Op] {
        p.op = Flt
    }

    return p
}

func Jmp(s *State, p pargs) {
    s.IP = p.jump
}

func JmpIfFlag(s *State, p pargs, flag uint, rev bool) {
    if ( (s.Reg[RF] & (1 << flag) > 0) == rev ) {
        Jmp(s, p)
    }
}

func Call(s *State, p pargs) {

    if s.pt == 31 {
        s.Reg[RF] |= (1 << F_FAULT)
        return
    }

    s.pt += 1
    s.Ret[s.pt] = s.IP
    s.IP = p.jump
}

func Ret(s *State, p pargs) {

    if s.pt == 0 {
        s.Reg[RF] |= (1 << F_FAULT)
        return
    }

    s.IP = s.Ret[s.pt]
    s.pt -= 1
}

func Jz(s *State, p pargs) {
    JmpIfFlag(s, p, F_ZERO, false)
}

func Jnz(s *State, p pargs) {
    JmpIfFlag(s, p, F_ZERO, true)
}

func Jo(s *State, p pargs) {
    JmpIfFlag(s, p, F_OVER, false)
}

func Jno(s *State, p pargs) {
    JmpIfFlag(s, p, F_OVER, true)
}

func Jf(s *State, p pargs) {
    JmpIfFlag(s, p, F_FAULT, false)
}

func Jnf(s *State, p pargs) {
    JmpIfFlag(s, p, F_FAULT, true)
}


func (s *State) result(result int) byte {
    if result < 0 {
        result = -result
        s.Reg[RF] |= (1 << F_FAULT)
    }

    if result > 0xFF {
        result = result & 0xFF
        s.Reg[RF] |= (1 << F_OVER)
    }


    if result == 0 {
        s.Reg[RF] |= (1 << F_ZERO)
    }

    return byte(result)
}


func Add(s *State, p pargs) {
    *p.R = s.result(int(p.A[0]) + int(p.A[1]))
}

func Sub(s *State, p pargs) {
    *p.R = s.result(int(p.A[0]) - int(p.A[1]))
}

func Inc(s *State, p pargs) {
    *p.R = s.result(int(p.A[0]) + 1)
}

func Dec(s *State, p pargs) {
    *p.R = s.result(int(p.A[0]) - 1)
}

func Mul(s *State, p pargs) {
    *p.R = s.result(int(p.A[0]) * int(p.A[1]))
}

func Div(s *State, p pargs) {
    if p.A[1] != 0 {
        *p.R = byte(int(p.A[0]) / int(p.A[1]))
        s.Reg[RA] = (byte) (int(p.A[0]) % int(p.A[1]))
    } else {
        s.Reg[RF] |= (1 << F_FAULT)
    }
}

func And(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] & p.A[1]))
}

func Or(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] | p.A[1]))
}

func Xor(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] ^ p.A[1]))
}

func Not(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] ^ 0xFF))
}

func Rol(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] << 1))
}

func Ror(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] >> 1))
}

func Rcl(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] << 1) | int(p.A[0] >> 7))
}

func Rcr(s *State, p pargs) {
    *p.R = s.result(int(p.A[0] >> 1) | int(p.A[0] << 7))
}

func Mov(s *State, p pargs) {
    *p.R = s.result(int(p.A[1]))
}

func Ld(s *State, p pargs) {
    *p.R = s.Mem[p.A[1]]
}

func St(s *State, p pargs) {
    s.Mem[p.A[1]] = p.A[0]
}

func Push(s *State, p pargs) {
    if s.Reg[RE] < 128 || s.Reg[RE] > 127 + 64 {
        s.Reg[RF] |= (1 << F_FAULT)
    } else {
        s.Mem[s.Reg[RE]] = p.A[0]
        s.Reg[RE] = byte(int(s.Reg[RE]) + 1)
    }
}

func Pop(s *State, p pargs) {
    if s.Reg[RE] < 128 || s.Reg[RE] > 127 + 64 {
        s.Reg[RF] |= (1 << F_FAULT)
    } else {
        s.Reg[RE] = byte(int(s.Reg[RE]) - 1)
        *p.R = s.result(int(s.Mem[s.Reg[RE]]))
    }
}


// Generate fault state
func Flt(s *State, p pargs) {
    s.Reg[RF] |= (1 << F_FAULT)
}

func Nop(s *State, p pargs) {
}
