package cpu

import (
    "fmt"

    "github.com/wthys/advent-of-code-2022/solutions/day10/cpu/register"
)

type (
    CPU struct {
        register *register.Register
        instr InstructionSet
        watchers []CPUWatcher
    }

    CPUWatcher func(instr Instruction, reg *register.Register)

    InstructionSet map[string]InstructionMaker

    InstructionMaker interface {
        Make(input string) (Instruction, error)
        Name() string
    }

    Instruction interface {
        Cost() int
        Execute(reg *register.Register)
    }
)


func New(makers ...InstructionMaker) *CPU {
    set := InstructionSet{}
    for _, maker := range makers {
        set[maker.Name()] = maker
    }
    reg := register.New()
    reg.Set("X", 1)
    reg.Set("cycle", 0)

    watchers := []CPUWatcher{}

    cpu := &CPU{reg, set, watchers}
    return cpu
}

func (cp *CPU) AddInstr(maker InstructionMaker) {
    cp.instr[maker.Name()] = maker
}

func (cp *CPU) Register() *register.Register {
    return cp.register
}

func (cp *CPU) cycle(instr Instruction) {
    cp.register.Inc("cycle", 1)
    for _, watcher := range cp.watchers {
        watcher(instr, cp.register)
    }
}

func (cp *CPU) Execute(program []Instruction) {
    for _, instr := range program {
        cost := instr.Cost()
        for i := 0; i < cost; i++ {
            cp.cycle(instr)
        }
        instr.Execute(cp.register)
    }
}

func (cp *CPU) AddWatcher(watcher CPUWatcher) {
    cp.watchers = append(cp.watchers, watcher)
}

func (cp CPU) Compile(input string) (Instruction, error) {
    for _, maker := range cp.instr {
        instr, err := maker.Make(input)
        if err != nil {
            return nil, err
        }

        if instr != nil {
            return instr, nil
        }
    }
    return nil, fmt.Errorf("Unknown instruction: %q", input)
}
