package mio

import (
    "os"
    "github.com/edsrzf/mmap-go"
)

type Prog struct {
    file *os.File
    mmap mmap.MMap
    Str []string
}

const MAXLEN = 2048

func (p *Prog) makedb() {

    p.Str = make([]string, 512, MAXLEN)

    var f, n int

    for c := 1; c < len(p.mmap) && c < MAXLEN; c += 1 {
        if p.mmap[c] == '\n' {
            p.Str[n] = string(p.mmap[f:c])
            n += 1
            c += 1
            f = c
        }
    }



}

func Mmap(path string) (p Prog, err error) {

    p.file, err = os.Open(path)

    if err != nil {
        return p, err
    }

    p.mmap, err = mmap.Map(p.file, mmap.RDONLY, 0)

    if err != nil {
        p.Close()
        return p, err
    }

    p.makedb()
    return p, nil
}

func (p *Prog) Close() {
    p.Str = nil
    p.mmap.Unmap()
    p.file.Close()
}


