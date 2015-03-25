package asm

import (
    . "github.com/derlaft/figex/mio"
    "testing"
    "fmt"
    //"strconv"
)

func TestTokenize(t *testing.T) {

    s := "\t\tMOV  AL \tAL "
    a := tokenize(s)


    if a[0] != "MOV" || a[1] != "AL" || a[2] != "AL" {
                t.Fail()
            }
}

func TestPreprocessJump(t *testing.T) {
    prog, err := ProgFromFile("./TEST1.PER")
        if err != nil {
            fmt.Println(err)
                t.Fail()
                return
        }

    a := AsmState{}
    Preprocess(prog.Str, &a)
    fmt.Printf("%q\n", a.Const)

    a.PC = 10
    fmt.Println(Cycle(prog, &a))
        fmt.Printf("%d %q\n", a.PC, a.Reg)
    fmt.Println(Cycle(prog, &a))
        fmt.Printf("%d %q\n", a.PC, a.Reg)
    if a.PC != 7 {
        t.Fail()
        return
    }
}




