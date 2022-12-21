package day21

import (
    "fmt"
    "strconv"
    "regexp"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/util"
    "github.com/wthys/advent-of-code-2022/collections"

)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "21"
}

func (s solution) Part1(input []string) (string, error) {
    monkeys, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    yelled := collections.NewSet[string]()

    notifyMonkeys := func(name string, number int) {
        for _, monkey := range monkeys {
            monkey.Listen(name, number)
        }
    }

    for !yelled.Has("root") {
        for name, monkey := range monkeys {
            if yelled.Has(name) {
                continue
            }

            number, err := monkey.Yell()
            if err != nil {
                continue
            }

            notifyMonkeys(name, number)
            yelled.Add(name)
        }

        //fmt.Printf("progress: yelled %4d/%4d\r", yelled.Len(), len(monkeys))
    }
    //fmt.Println()

    number, err := monkeys["root"].Yell()
    if err != nil {
        return solver.Error(err)
    }
    return solver.Solved(number)
}

func (s solution) Part2(input []string) (string, error) {
    monkeys, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    order := map[string]int{}
    invorder := map[int][]string{}
    yellers := map[string]Monkey{}

    addOrder := func(name string, n int) {
        names, ok := invorder[n]
        if !ok {
            names = []string{}
        }
        invorder[n] = append(names, name)

        order[name] = n
    }

    determined := collections.NewSet[string]()
    maxOrder := 0

    for determined.Len() < len(monkeys) {
        for name, monkey := range monkeys {
            if determined.Has(name) {
                continue
            }

            _, ok := monkey.(*YellerMonkey)
            if ok {
                yellers[name] = monkey
                addOrder(name, 0)
                determined.Add(name)
                continue
            }

            math := *(monkey.(*MathMonkey))

            if !(determined.Has(math.Left) && determined.Has(math.Right)) {
                continue
            }

            addOrder(name, util.Max(order[math.Left], order[math.Right]) + 1)
            if order[name] > maxOrder {
                maxOrder = order[name]
            }
            determined.Add(name)
        }
    }

    (*(monkeys["root"].(*MathMonkey))).Operation = "-"


    notifyMonkeys := func(name string, number int) {
        for _, monkey := range monkeys {
            monkey.Listen(name, number)
        }
    }
    excersizeMonkeys := func(number int) (int, error) {
        monkeys["humn"] = &YellerMonkey{number}

        for _, monkey := range monkeys {
            monkey.Reset()
        }

        for i := 0; i <= maxOrder; i++ {
            for _, name := range invorder[i] {
                number, err := monkeys[name].Yell()
                if err != nil {
                    return 0, err
                }
                notifyMonkeys(name, number)
            }
        }

        return monkeys["root"].Yell()
    }

    humanShout := BinarySearch(func(number int) int {
        val, _ := excersizeMonkeys(number)
        return val
    })

    zeroes := []int{}

    for i:= -6; i <= 6; i++ {
        val, _ := excersizeMonkeys(humanShout + i)
        if val == 0 {
            zeroes = append(zeroes, humanShout + i)
        }
        //fmt.Printf("CHECK => when shouting %v, root shouts %v\n", humanShout+i, val)
    }

    return solver.Solved(util.Min(zeroes...))
}

func BinarySearch(fn func(number int) int) int {
    lo := 0
    hi := 1

    ylo := fn(lo)
    yhi := fn(hi)

    for util.Sign(ylo) == util.Sign(yhi) {
        lo, ylo = hi, yhi
        hi = hi * 2
        yhi = fn(hi)
    }

    for {
        mid := (lo + hi) / 2
        ymid := fn(mid)

        if ymid == 0 {
            return mid
        }

        if util.Sign(ylo) != util.Sign(ymid) {
            hi = mid
            yhi = ymid
        } else {
            lo = mid
            ylo = ymid
        }
    }
}

type (
    KnownNumbers map[string]int

    Monkey interface {
        Yell() (int, error)
        Listen(name string, number int)
        Reset()
    }

    YellerMonkey struct {
        Number int
    }

    MathMonkey struct {
        Left string
        Right string
        Operation string
        Overheard KnownNumbers
    }

    MathOp func(a, b int) int
)

var (
    MATHOPS = map[string]MathOp{
        "*": func(a, b int) int { return a * b },
        "+": func(a, b int) int { return a + b },
        "/": func(a, b int) int { return a / b },
        "-": func(a, b int) int { return a - b },
    }
)

func (m YellerMonkey) Yell() (int, error) {
    return m.Number, nil
}

func (m YellerMonkey) Listen(_ string, _ int) {
    return
}

func (m YellerMonkey) Reset() { return }

func (m YellerMonkey) String() string {
    return fmt.Sprintf("YellerMonkey{%v}", m.Number)
}

func (m MathMonkey) Yell() (int, error) {
    if m.IsWaiting() {
        return 0, fmt.Errorf("have not heard right monkeys")
    }

    a := m.Overheard[m.Left]
    b := m.Overheard[m.Right]

    return MATHOPS[m.Operation](a, b), nil
}

func (m *MathMonkey) Listen(name string, number int) {
    if name != (*m).Left && name != (*m).Right {
        return
    }

    (*m).Overheard[name] = number
}

func (m MathMonkey) IsWaiting() bool {
    _, present := m.Overheard[m.Left]
    if !present {
        return true
    }

    _, present = m.Overheard[m.Right]
    return !present
}

func (m *MathMonkey) Reset() {
    (*m).Overheard = KnownNumbers{}
}

func (m MathMonkey) String() string {
    return fmt.Sprintf("MathMonkey{%v %v %v}", m.Left, m.Operation, m.Right)
}

func parseInput(input []string) (map[string]Monkey, error) {
    reMath := regexp.MustCompile("^([a-z]+):\\s*([a-z]+)\\s+([*/+-])\\s+([a-z]+)\\s*$")
    reYell := regexp.MustCompile("^([a-z]+):\\s*([0-9]+)\\s*$")

    monkeys := map[string]Monkey{}

    for nr, line := range input {
        caps := reMath.FindStringSubmatch(line)
        if caps != nil {
            monkeys[caps[1]] = &MathMonkey{caps[2], caps[4], caps[3], KnownNumbers{}}
            continue
        }

        caps = reYell.FindStringSubmatch(line)
        if caps != nil {
            num, err := strconv.Atoi(caps[2])
            if err != nil {
                return nil, err
            }

            monkeys[caps[1]] = &YellerMonkey{num}
            continue
        }
        return nil, fmt.Errorf("could not understand line #%v : %v", nr, line)
    }

    return monkeys, nil
}
