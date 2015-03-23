package cpu

import (
    "cpu"
)

type AsmState struct {
    cpu.State
    Labels map[string]int
}
