package register

import (
    "fmt"
)

type Register struct {
    regs map[string]int
}

func New() *Register {
    return &Register{map[string]int{}}
}

func (r Register) String() string {
    return fmt.Sprintf("%v", r.regs)
}

func (r Register) Get(name string) int {
    val, ok := r.regs[name]
    if !ok {
        val = 0
    }
    return val
}

func (r *Register) Set(name string, value int) {
    r.regs[name] = value
}

func (r *Register) Inc(name string, value int) {
    r.regs[name] = r.Get(name) + value
}
