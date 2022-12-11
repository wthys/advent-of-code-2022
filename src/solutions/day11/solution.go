package day11

import (
    "fmt"
    "sort"
    "regexp"
    "strconv"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "11"
}

func (s solution) Part1(input []string) (string, error) {
    monkeys, err := parseInput(input)
    if err != nil {
        return "", err
    }


    // Setup listener
    inspections := map[string]int{}
    for key, monkey := range monkeys {
        k := key
        //fmt.Printf("Monkey %v: %v\n", key, monkey.Inventory())
        inspections[k] = 0
        lsnr := func() {
            inspections[k] += 1
        }
        monkey.AddListener(lsnr)
    }

    rounds := 20

    for i := 0; i < rounds; i++ {
        for id := 0; id < len(monkeys); id++ {
            key := fmt.Sprint(id)
            //fmt.Printf("Monkey %v:\n", key)
            monkeys[key].InspectAndThrow(monkeys)
        }
    }

    values := sort.IntSlice{}
    for _, insp := range inspections {
        //fmt.Printf("Monkey %v inspected %v times\n", key, insp)
        values = append(values, insp)
    }
    sort.Sort(sort.Reverse(values))

    business := values[0] * values[1]

    return fmt.Sprint(business), nil
}

func (s solution) Part2(input []string) (string, error) {
    monkeys, err := parseInput(input)
    if err != nil {
        return "", err
    }

    // Setup listener
    inspections := map[string]int{}
    lcd := 1
    for _, monkey := range monkeys {
        lcd *= monkey.test
    }

    relief := func(v int) int {
        return v % lcd
    }

    for key, monkey := range monkeys {
        k := key
        //fmt.Printf("Monkey %v: %v\n", key, monkey.Inventory())
        inspections[k] = 0
        lsnr := func() {
            inspections[k] += 1
        }
        monkey.AddListener(lsnr)
        monkey.relief = relief
    }

    rounds := 10000

    for i := 1; i <= rounds; i++ {
        for id := 0; id < len(monkeys); id++ {
            key := fmt.Sprint(id)
            //fmt.Printf("Monkey %v:\n", key)
            monkeys[key].InspectAndThrow(monkeys)
        }

        /*
        if i == 1 || i == 20 || (i > 0 && i % 1000 == 0) {
            fmt.Printf("== After round %v ==\n", i)
            for key, insp := range inspections {
                fmt.Printf("Monkey %v inspected %v times\n", key, insp)
            }
            fmt.Println()
        } else {
            fmt.Printf(":: checking round %v     \r", i)
        }
        //*/
    }

    values := sort.IntSlice{}
    for _, insp := range inspections {
        //fmt.Printf("Monkey %v inspected %v times\n", key, insp)
        values = append(values, insp)
    }
    sort.Sort(sort.Reverse(values))

    business := values[0] * values[1]

    return fmt.Sprint(business), nil
}

type (
    Monkey struct {
        items []int
        oper Operation
        test int
        action map[bool]string
        relief ReliefFunction
        listeners []ListenerFunction
    }

    ReliefFunction func(int) int
    ListenerFunction func()

    Operation interface {
        NewValue(old int) int
    }

    MultOper struct {
        factor int
    }

    AddOper struct {
        term int
    }

    SquareOper struct {}
)

func (op SquareOper) NewValue(old int) int {
    return old * old
}

func (op AddOper) NewValue(old int) int {
    return old + op.term
}

func (op MultOper) NewValue(old int) int {
    return old * op.factor
}

func (m *Monkey) Inventory() []int {
    return m.items
}

func (m *Monkey) InspectAndThrow(monkeys map[string]*Monkey) {
    //oknok := map[bool]string{true: " ", false: " not "}
    for _, worry := range m.items {
        //fmt.Printf("  Monkey inspects an item with a worry level of %v\n", worry)
        worry = m.oper.NewValue(worry)
        //fmt.Printf("    Worry level is adjusted to %v\n", worry)
        worry = m.relief(worry)
        //fmt.Printf("    Worry level divided by 3 to %v\n", worry)
        for _, listener := range m.listeners {
            listener()
        }
        test := worry % m.test == 0
        target := m.action[test]
        //fmt.Printf("    Worry level is%vdivisible by %v, throwing to monkey %v\n", oknok[test], m.test, target)
        monkeys[target].Catch(worry)
    }

    m.items = []int{}
}

func (m *Monkey) Catch(worry int) {
    if m.items == nil {
        m.items = []int{}
    }
    m.items = append(m.items, worry)
}

func (m *Monkey) AddListener(listener ListenerFunction) {
    m.listeners = append(m.listeners, listener)
}

func parseInput(input []string) (map[string]*Monkey, error) {
    reMonkey := regexp.MustCompile("^Monkey (\\d+):\\s*$")
    reStart := regexp.MustCompile("^  Starting items: (\\d+(, \\d+)*)\\s*$")
    reOper := regexp.MustCompile("^  Operation: new = old ([*+]) (old|\\d+)\\s*$")
    reTest := regexp.MustCompile("^  Test: divisible by (\\d+)\\s*$")
    reCond := regexp.MustCompile("^    If (true|false): throw to monkey (\\d+)\\s*$")

    div3 := func(v int) int {
        return v / 3
    }

    monkeys := map[string]*Monkey{}
    for i := 0; i < len(input); i += 7 {
        lines := input[i:i+6]

        caps := reMonkey.FindStringSubmatch(lines[0])
        if caps == nil {
            return nil, fmt.Errorf("Could not parse Monkey %v line 0", i)
        }
        key := caps[1]

        caps = reStart.FindStringSubmatch(lines[1])
        if caps == nil {
            return nil, fmt.Errorf("Could not parse Monkey %v line 1", i)
        }
        items := []int{}
        for _, item := range strings.Split(caps[1], ",") {
            val, _ := strconv.Atoi(strings.TrimSpace(item))
            items = append(items, val)
        }

        caps = reOper.FindStringSubmatch(lines[2])
        if caps == nil {
            return nil, fmt.Errorf("Could not parse Monkey %v line 2", i)
        }
        var oper Operation
        switch caps[1] {
        case "+":
            if caps[2] == "old" {
                oper = MultOper{2}
            } else {
                val, _ := strconv.Atoi(caps[2])
                oper = AddOper{val}
            }
        default:
            if caps[2] == "old" {
                oper = SquareOper{}
            } else {
                val, _ := strconv.Atoi(caps[2])
                oper = MultOper{val}
            }
        }

        caps = reTest.FindStringSubmatch(lines[3])
        if caps == nil {
            return nil, fmt.Errorf("Could not parse Monkey %v line 3", i)
        }
        test, _ := strconv.Atoi(caps[1])

        caps = reCond.FindStringSubmatch(lines[4])
        if caps == nil {
            return nil, fmt.Errorf("Could not parse Monkey %v line 4", i)
        }
        yes := caps[2]

        caps = reCond.FindStringSubmatch(lines[5])
        if caps == nil {
            return nil, fmt.Errorf("Could not parse Monkey %v line 5", i)
        }
        no := caps[2]

        monkeys[key] = &Monkey{items, oper, test, map[bool]string{true:yes, false:no}, div3, []ListenerFunction{}}
    }

    return monkeys, nil
}
