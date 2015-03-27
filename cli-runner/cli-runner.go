package main

import (
    . "github.com/derlaft/figex/mio"
    . "github.com/derlaft/figex/asm"
    "os"
    "fmt"
    "bufio"
)


func main() {
    fname := os.Args[1]

    prog, err := ProgFromFile(fname)

    if err != nil {
        fmt.Println(err)
    }

    state := State{}

    scanner := bufio.NewScanner(os.Stdin)

    for i, in := range prog.Prog {
        fmt.Println(i, in.Op)
    }

    for scanner.Scan() && state.GetIP() < len(prog.Prog) {
        args := prog.Prog[state.GetIP()]
        state.Tick(args)
        fmt.Printf("%q\n", state.Reg)
        fmt.Printf("%s\n", args.Op)
    }

}


