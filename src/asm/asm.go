package cpu

import (
    . "cpu"
    "strings"
    "strconv"
    "errors"
)

type AsmState struct {
    State
    PC int
    Const map[string]int
}

type OrgexPut struct {
    Orgex
    used int
}

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
}

// instruction that need Orgex.R pointer to return value
var Returning = map[string]bool {
    "ADD": true,
    "SUB": true,
    "INC": true,
    "DEC": true,
    "MUL": true,
    "DIV": true,
    "AND": true,
    "OR": true,
    "XOR": true,
    "NOT": true,
    "ROL": true,
    "ROR": true,
    "RCL": true,
    "RCR": true,
    "MOV": true,
    "LD": true,
    "POP": true,
}

// @TODO: check for interrupts
// and implement them ofc :)

func Cycle(a []string, state AsmState) error {

    var token []string

    for argType(token) != OP_OP && state.PC < len(a) {
        token = tokenize(a[state.PC])
        state.PC += 1
    }

    if state.PC >= len(a) {
        return errors.New("Progam ended")
    }

    op := token[0]
    args := OrgexPut{}

    for _, arg := range token[1:] {
        err := args.pushArg(state, arg)
        if err != nil {
            return errors.New("Parse failed on line " + strconv.Itoa(state.PC))
        }
    }

    _, has := Ops[op]
    if !has {
        return errors.New("Unknown instruction")
    }

    _, has = Returning[op]
    if !has && args.R == nil {
        state.Reg[RF] |= (1 << F_FAULT)
        //@TODO: destroy stack pointer
        return nil
    }

    state.Reg[RF] = 0
    Ops[op](&state.State, args.Orgex)
    state.PC += 1


    return nil
}

func (org *OrgexPut) pushArg(state AsmState, arg string) error {

    var res byte

    number, err := strconv.ParseInt(arg[1:], 16, 8)
    if err != nil {
        return err
    }


    switch arg[0] {
        //hex constant
        case '&':
            res = byte(number)
        case 'R':
            if number >= 16 {
                return errors.New("")
            }
            res = state.Reg[number]
            if org.used == 0 {
                org.R = &state.Reg[number]
            }
    }

    switch org.used {
        case 0:
            org.A = res
        case 1:
            org.B = res
    }

    org.used += 1

    return nil

}

var (
    OP_OP = 0
    OP_LABEL = 1
    OP_CONSTANT = 2
    OP_NOP = -1
)

func tokenize(str string) []string {
    return strings.Split(strings.Trim(str, " \t"), " \t")
}

func argType(t []string) int {
    first, last := t[0][0], t[0][len(t[0])-1]
    switch {
        case first == '#' && t[0] == "#DEF" && len(t) == 3:
            return OP_CONSTANT
        case last == ':' && len(t) == 2:
            return OP_LABEL
        case first == '%' || len(t[0]) == 0:
            return OP_NOP
        default:
            return OP_OP
    }
}

func Preprocess(t []string, state *AsmState) {
    state.Const = make(map[string]int)
    for _, s := range t {
        tokens := tokenize(s)
        switch argType(tokens) {

            //@TODO: check conv error

            case OP_LABEL:
                name := tokens[0][0:len(tokens[0])-1]
                value, _ := strconv.Atoi(tokens[1])
                state.Const[name] = value

            case OP_CONSTANT:
                value, _ := strconv.Atoi(tokens[2])
                state.Const[tokens[1]] = value
        }
    }
}


