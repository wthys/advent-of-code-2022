package day10

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/solutions/day10/cpu"
    "github.com/wthys/advent-of-code-2022/solutions/day10/cpu/register"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "10"
}

func (s solution) Part1(input []string) (string, error) {
    cp := InitCpu()

    program, err := parseInput(cp, input)

    if err != nil {
        return "", err
    }

    signal := 0

    watcher := func(_ cpu.Instruction, reg *register.Register) {
        cycle := reg.Get("cycle")
        if (cycle-20)%40 != 0 {
            return
        }
        x := reg.Get("X")

        str := x * cycle

        signal += str
    }

    cp.AddWatcher(watcher)

    cp.Execute(program)

    return strconv.Itoa(signal), nil
}

func (s solution) Part2(input []string) (string, error) {
    cp := InitCpu()

    program, err := parseInput(cp, input)

    if err != nil {
        return "", err
    }

    crt := [][]rune{}
    for j := 0; j < 6; j++ {
        crt = append(crt, []rune(strings.Repeat("_", 40)))
    }

    watcher := func(_ cpu.Instruction, reg *register.Register) {
        cycle := reg.Get("cycle")
        x := reg.Get("X")

        row := (cycle - 1) / 40
        col := (cycle - 1) % 40

        if x == col-1 || x == col || x == col+1 {
            crt[row][col] = '#'
        }
    }

    cp.AddWatcher(watcher)

    cp.Execute(program)

    for _, line := range crt {
        fmt.Println(string(line))
    }

    return "^^^", nil
}


func InitCpu() *cpu.CPU {
    noop := NoopInstrMaker{nil, nil}
    addx := AddXInstrMaker{nil}
    return cpu.New(&noop, &addx)
}

func parseInput(c *cpu.CPU, input []string) ([]cpu.Instruction, error) {
    program := []cpu.Instruction{}

    for n, line := range input {
        instr, err := c.Compile(line)
        if err != nil {
            return nil, fmt.Errorf("Error on line %v : %v", n, err)
        }
        program = append(program, instr)
    }

    return program, nil
}


type (
    NoopInstr struct {}
    NoopInstrMaker struct{
        re *regexp.Regexp
        instr *NoopInstr
    }

    AddXInstr struct {
        amount int
    }
    AddXInstrMaker struct{
        re *regexp.Regexp
    }
)


func (instr *AddXInstr) Cost() int {
    return 2
}

func (instr *AddXInstr) Execute(reg *register.Register) {
    reg.Inc("X", instr.amount)
}


func (instr *NoopInstr) Cost() int {
    return 1
}

func (instr *NoopInstr) Execute(_ *register.Register) {
    return
}


func (m *AddXInstrMaker) Name() string {
    return "addx"
}

func (m *AddXInstrMaker) Make(input string) (cpu.Instruction, error) {
    if m.re == nil {
        m.re = regexp.MustCompile("^\\s*addx(\\s+(-?\\d+))?\\s*$")
    }

    caps := m.re.FindStringSubmatch(input)
    if caps == nil {
        return nil, nil
    }

    if len(caps) <= 1 {
        return nil, fmt.Errorf("Missing argument: %q", input)
    }

    amount, err := strconv.Atoi(caps[2])
    if err != nil {
        return nil, err
    }

    return &AddXInstr{amount}, nil
}


func (m *NoopInstrMaker) Name() string {
    return "noop"
}

func (m *NoopInstrMaker) Make(input string) (cpu.Instruction, error) {
    if m.re == nil {
        m.re = regexp.MustCompile("^\\s*noop\\s*$")
    }

    caps := m.re.FindStringSubmatch(input)
    if caps == nil {
        return nil, nil
    }

    if m.instr == nil {
        m.instr = &NoopInstr{}
    }

    return m.instr, nil
}
