package day5

import (
    "fmt"
    "strconv"
    "regexp"
    "strings"
    "github.com/golang-collections/collections/stack"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "5"
}


func tryStack(stacks map[int]*stack.Stack, idx int, name string) map[int]*stack.Stack {

    st, ok := stacks[idx]
    if !ok {
        st = stack.New()
        stacks[idx] = st
    }

    if name != "" {
        st.Push(name)
    }

    return stacks
}


func parseStacks(input []string) map[int]*stack.Stack {
    stacks := map[int]*stack.Stack{}

    level := regexp.MustCompile("(....)|(...$)")
    crate := regexp.MustCompile("[A-Z]")

    for _, line := range input {
        crates := level.FindAllString(line, -1)

        for i, c := range crates {
            name := crate.FindString(c)
            stacks = tryStack(stacks, i+1, name)
        }
    }

    proper := map[int]*stack.Stack{}

    for i, s := range stacks {
        p := stack.New()
        for s.Peek() != nil {
            p.Push(s.Pop())
        }
        proper[i] = p
    }

    return proper
}

type Instr struct {
    Amount int
    From int
    To int
}

func (i Instr) String() string {
    return fmt.Sprintf("move %v from %v to %v", i.Amount, i.From, i.To)
}

func parseInstr(input []string) []Instr {
    instr := make([]Instr, 0)

    num := regexp.MustCompile("[0-9]+")

    for n, line := range input {
        pos := num.FindAllString(line, 3)

        if len(pos) == 0 {
            fmt.Printf("%v : %q\n", n, line)
        }

        amount, _ := strconv.Atoi(pos[0])
        from, _ := strconv.Atoi(pos[1])
        to, _ := strconv.Atoi(pos[2])

        instr = append(instr, Instr{amount, from, to})
    }

    return instr
}

func parseInput(input []string) (map[int]*stack.Stack, []Instr) {
    sep := 0
    for i, line := range input {
        if len(strings.TrimSpace(line)) == 0 {
            sep = i
            break
        }
    }

    stacks := parseStacks(input[:sep])
    instr := parseInstr(input[sep+1:])

    return stacks, instr
}

func getMessage(stacks map[int]*stack.Stack) string {
    b := strings.Builder{}
    for i := 1; i <= len(stacks); i++ {
        b.WriteString(stacks[i].Peek().(string))
    }
    return b.String()
}


func (s solution) Part1(input []string) (string, error) {
    stacks, instr := parseInput(input)

    for _, in := range instr {
        for i := 0; i < in.Amount; i++ {
            stacks[in.To].Push(stacks[in.From].Pop())
        }
    }

    return getMessage(stacks), nil
}

func (s solution) Part2(input []string) (string, error) {
    stacks, instr := parseInput(input)

    for _, in := range instr {
        crane := stack.New()
        for i := 0; i < in.Amount; i++ {
            crane.Push(stacks[in.From].Pop())
        }
        for i := 0; i < in.Amount; i++ {
            stacks[in.To].Push(crane.Pop())
        }
    }

    return getMessage(stacks), nil
}
