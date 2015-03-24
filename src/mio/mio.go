package mio

import (
    "os"
    "bufio"
)

type Prog struct {
    Str []string
}

const MAXLEN = 2048

func Mmap(path string) (p Prog, err error) {
    p.Str = make([]string, 512, MAXLEN)

    file, err := os.Open(path)

    if err != nil {
        return p, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    n := 0
    for ; scanner.Scan(); n++ {
            p.Str[n] = scanner.Text()
    }

    //cut slice so probably GC will free some space
    p.Str = p.Str[:n]

    return p, nil
}

