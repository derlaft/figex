package asm

import (
    . "mio"
    "testing"
    "fmt"
    "strconv"
)

func TestTokenize(t *testing.T) {

    s := "\t\tMOV  AL \tAL "
    a := tokenize(s)


    if a[0] != "MOV" ||
       a[1] != "AL" ||
       a[2] != "AL" {
                t.Fail()
            }
}

func TestPreprocess(t *testing.T) {
    mmap, err := Mmap("./TEST1.PER")
        if err != nil {
            fmt.Println(err)
            t.Fail()
        }

    a := AsmState{}
    Preprocess(mmap.Str, &a)

    fmt.Println("Printing const table")
    for s, i := range a.Const {
        fmt.Println("Const " + s + " is " + strconv.Itoa(i))
    }


    defer mmap.Close()

}

