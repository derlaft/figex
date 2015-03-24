package io

import (
    "testing"
    "fmt"
)

func TestOpen(t *testing.T) {

    mmap, err := Mmap("./test.txt")
    if err != nil {
        fmt.Println(err)
        t.Fail()
    }

    if  mmap.Str[0] != "TEST1" ||
        mmap.Str[1] != "TEST2" ||
        mmap.Str[2] != "MEOWOWOWOWWOWOOW?" ||
        mmap.Str[3] != "PUR PUR ;MAGIC" ||
        mmap.Str[4] != "Nyan)" {
            t.Fail()
        }

    defer mmap.Close()
}

