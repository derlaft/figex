package cpu

type State struct {
    Reg [16]byte
    Mem [256]byte
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

type Orgex struct {
    A byte
    B byte
    R *byte
}


type Instruction func(*State,Orgex)

func (s *State) result(result int) byte {
    if result > 0xFF {
        result = result & 0xFF
        s.Reg[RF] |= (1 << F_OVER)
    }

    if result < 0 {
        result = -result
        s.Reg[RF] |= (1 << F_FAULT)
    }

    if result == 0 {
        s.Reg[RF] |= (1 << F_ZERO)
    }

    return byte(result)
}


func Add(s *State, p Orgex) {
    *p.R = s.result(int(p.A) + int(p.B))
}

func Sub(s *State, p Orgex) {
    *p.R = s.result(int(p.A) - int(p.B))
}

func Inc(s *State, p Orgex) {
    *p.R = s.result(int(p.A) + 1)
}

func Dec(s *State, p Orgex) {
    *p.R = s.result(int(p.A) - 1)
}

func Mul(s *State, p Orgex) {
    *p.R = s.result(int(p.A) * int(p.B))
}

func Div(s *State, p Orgex) {
    if p.B != 0 {
        *p.R = byte(int(p.A) / int(p.B))
        s.Reg[RA] = (byte) (int(p.A) % int(p.B))
    } else {
        s.Reg[RF] |= (1 << F_FAULT)
    }
}

func And(s *State, p Orgex) {
    *p.R = s.result(int(p.A & p.B))
}

func Or(s *State, p Orgex) {
    *p.R = s.result(int(p.A | p.B))
}

func Xor(s *State, p Orgex) {
    *p.R = s.result(int(p.A ^ p.B))
}

func Not(s *State, p Orgex) {
    *p.R = s.result(int(p.A ^ 0xFF))
}

func Rol(s *State, p Orgex) {
    *p.R = s.result(int(p.A << 1))
}

func Ror(s *State, p Orgex) {
    *p.R = s.result(int(p.A >> 1))
}

func Rсl(s *State, p Orgex) {
    *p.R = s.result(int(p.A << 1) | int(p.A >> 7))
}

func Rcr(s *State, p Orgex) {
    *p.R = s.result(int(p.A >> 1) | int(p.A << 7))
}

func Mov(s *State, p Orgex) {
    *p.R = p.B
}

func Ld(s *State, p Orgex) {
    *p.R = s.Mem[p.B]
}

func St(s *State, p Orgex) {
    s.Mem[p.B] = p.A
}

func Push(s *State, p Orgex) {
    if s.Reg[RE] < 128 || s.Reg[RE] > 127 + 64 {
        s.Reg[RF] |= (1 << F_FAULT)
    } else {
        s.Mem[s.Reg[RE]] = p.A
        s.Reg[RE] = byte(int(s.Reg[RE]) + 1)
    }
}

func Pop(s *State, p Orgex) {
    if s.Reg[RE] < 128 || s.Reg[RE] > 127 + 64 {
        s.Reg[RF] |= (1 << F_FAULT)
    } else {
        s.Reg[RE] = byte(int(s.Reg[RE]) - 1)
        *p.R = s.result(int(s.Mem[s.Reg[RE]]))
    }
}

