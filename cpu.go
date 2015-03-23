package figex

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
  RE = 14
  RF = 15

  F_ZERO = 0
  F_OVER = 1
  F_FAULT = 2
  F_INT = 3
)

type Pair struct {
    A byte
    B byte
    used byte
}

type instruction func(State,Pair)

func (s *State) result(result int) byte {
    if result > 0xFF {
        result = result % 0xFF
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


func add(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A]) + int(s.Reg[par.B]))
}

func sub(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A]) - int(s.Reg[par.B]))
}

func inc(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A]) + 1)
}

func dec(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A]) - 1)
}

func mul(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A]) * int(s.Reg[par.B]))
}

func div(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A]) / int(s.Reg[par.B]))
    s.Reg[RA] = (byte) (int(s.Reg[par.A]) % int(s.Reg[par.B]))
}

func and(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A] & s.Reg[par.B]))
}

func or(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A] | s.Reg[par.B]))
}

func xor(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A] ^ s.Reg[par.B]))
}

func not(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A] ^ 0xFF))
}

func rol(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A] << 1))
}

func ror(s State, par Pair) {
    s.Reg[par.A] = s.result(int(s.Reg[par.A] >> 1))
}

func mov(s State, par Pair) {
    s.Reg[par.A] = s.Reg[par.B]
}

func ld(s State, par Pair) {
    s.Reg[par.A] = s.Mem[s.Reg[par.B]]
}

func st(s State, par Pair) {
    s.Mem[s.Reg[par.B]] = s.Reg[par.A]
}




