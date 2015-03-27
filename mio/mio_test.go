package mio

import (
    "testing"
    "fmt"
)

func TestOpen(t *testing.T) {

    prog, err := ProgFromFile("./test.txt")
    if err != nil {
        fmt.Println(err)
        t.Fail()
    }

    if  prog[0] != "TEST1" ||
        prog[1] != "TEST2" ||
        prog[2] != "MEOWOWOWOWWOWOOW?" ||
        prog[3] != "PUR PUR ;MAGIC" ||
        prog[4] != "Nyan)" {
            t.Fail()
        }
}

func TestTokenize(t *testing.T) {

    s := "\t\tMOV  AL \tAL "
    a := tokenize(s)


    if a[0] != "MOV" || a[1] != "AL" || a[2] != "AL" {
        t.Fail()
    }
}

