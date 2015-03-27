package asm

import (
    "testing"
    "fmt"
    "strconv"
)

func TestResult(t *testing.T) {
    var good = map[int]byte {
        42:     0,
        15:     0,
        250:    0,
        -42:    (1 << F_FAULT),
        3556:   (1 << F_OVER),
        -4123:  (1 << F_OVER) | (1 << F_FAULT),
        0:      (1 << F_ZERO),
        256:    (1 << F_ZERO) | (1 << F_OVER),
    }

    for k, v := range good {
        state := State{}
        state.result(k)
        res := state.Reg[RF]
        if res != v {
            fmt.Println("Test " + strconv.Itoa(k))
            fmt.Println(strconv.Itoa(int(res)) + " != " + strconv.Itoa(int(v)))
            t.Fail()
        }
    }
}

