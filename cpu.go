package cpu

type State struct {
    Reg [16]byte
    Mem [256]byte
    PC int
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
    used byte
}

type instruction func(State,Orgex)

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


func add(s *State, p Orgex) {
    *p.R = s.result(int(p.A) + int(p.B))
}

func sub(s *State, p Orgex) {
    *p.R = s.result(int(p.A) - int(p.B))
}

func inc(s *State, p Orgex) {
    *p.R = s.result(int(p.A) + 1)
}

func dec(s *State, p Orgex) {
    *p.R = s.result(int(p.A) - 1)
}

func mul(s *State, p Orgex) {
    *p.R = s.result(int(p.A) * int(p.B))
}

func div(s *State, p Orgex) {
    if p.B != 0 {
        *p.R = byte(int(p.A) / int(p.B))
        s.Reg[RA] = (byte) (int(p.A) % int(p.B))
    } else {
        s.Reg[RF] |= (1 << F_FAULT)
    }
}

func and(s *State, p Orgex) {
    *p.R = s.result(int(p.A & p.B))
}

func or(s *State, p Orgex) {
    *p.R = s.result(int(p.A | p.B))
}

func xor(s *State, p Orgex) {
    *p.R = s.result(int(p.A ^ p.B))
}

func not(s *State, p Orgex) {
    *p.R = s.result(int(p.A ^ 0xFF))
}

func rol(s *State, p Orgex) {
    *p.R = s.result(int(p.A << 1))
}

func ror(s *State, p Orgex) {
    *p.R = s.result(int(p.A >> 1))
}

func mov(s *State, p Orgex) {
    *p.R = p.B
}

func ld(s *State, p Orgex) {
    *p.R = s.Mem[p.B]
}

func st(s *State, p Orgex) {
    s.Mem[p.B] = p.A
}

func push(s *State, p Orgex) {
    if s.Reg[RE] < 128 || s.Reg[RE] > 128 + 64 {
        s.Reg[RF] |= (1 << F_FAULT)
    } else {
        s.Mem[s.Reg[RE]] = p.A
        s.Reg[RE] = byte(int(s.Reg[RE]) + 1)
    }
}

func pop(s *State, p Orgex) {
    if s.Reg[RE] < 128 || s.Reg[RE] > 127 + 64 {
        s.Reg[RF] |= (1 << F_FAULT)
    } else {
        s.Reg[RE] = byte(int(s.Reg[RE]) - 1)
        *p.R = s.result(int(s.Mem[s.Reg[RE]]))
    }
}






