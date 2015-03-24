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

    if  prog.Str[0] != "TEST1" ||
        prog.Str[1] != "TEST2" ||
        prog.Str[2] != "MEOWOWOWOWWOWOOW?" ||
        prog.Str[3] != "PUR PUR ;MAGIC" ||
        prog.Str[4] != "Nyan)" {
            t.Fail()
        }

}

