package day20

import (
    "fmt"
    "strconv"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/util"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "20"
}

func (s solution) Part1(input []string) (string, error) {
    items, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    mixed := MixList(items)

    zeroIdx := 0
    for i, item := range mixed {
        if item.Value == 0 {
            zeroIdx = i
            break
        }
    }

    ll := len(mixed)
    k1 := zeroIdx + 1000
    k2 := zeroIdx + 2000
    k3 := zeroIdx + 3000

    a := mixed[k1 % ll].Value
    b := mixed[k2 % ll].Value
    c := mixed[k3 % ll].Value
    sum := a + b + c
    //fmt.Printf("%v + %v + %v = %v\n", a, b, c, sum)

    return solver.Solved(sum)
}

func (s solution) Part2(input []string) (string, error) {
    items, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }
    factor := 811589153

    mixed := []Item{}
    for _, item := range items {
        mixed = append(mixed, Item{item.Value * factor, item.Index})
    }

    //fmt.Printf("initial : %v\n", mixed)
    for n := 0; n < 10; n++ {
        mixed = MixList(mixed)
        //fmt.Printf("round %2d: %v\n", n+1, mixed)
    }


    zeroIdx := 0
    for i, item := range mixed {
        if item.Value == 0 {
            zeroIdx = i
            break
        }
    }

    ll := len(mixed)
    k1 := zeroIdx + 1000
    k2 := zeroIdx + 2000
    k3 := zeroIdx + 3000

    a := mixed[k1 % ll].Value
    b := mixed[k2 % ll].Value
    c := mixed[k3 % ll].Value
    sum := a + b + c
    //fmt.Printf("%v + %v + %v = %v\n", a, b, c, sum)

    return solver.Solved(sum)
}

type (
    Item struct {
        Value int
        Index int
    }
)

func (i Item) String() string {
    return fmt.Sprintf("%v:%v", i.Value, i.Index)
}

func MixList(items []Item) []Item {
    mixed := items
    for i, _ := range items {
        mixed = Mix(mixed, i)
    }
    return mixed
}

func Mix(items []Item, index int) []Item {
    idx := 0
    toMix := Item{}

    for i, item := range items {
        if item.Index == index {
            idx = i
            toMix = item
            break
        }
    }

    length := len(items)

    dir := util.Sign(toMix.Value)
    dist := util.Abs(toMix.Value) % (length * (length - 1))

    //shift other values to the left
    shift := (dist / length) % (length - 1)
    result := items
    if shift != 0 {
        others := []Item{}
        if idx >= length-1 {
            others = result[:length-1]
        } else if idx == 0 {
            others = result[1:]
        } else {
            others = append(result[:idx], result[idx+1:]...)
        }

        if dir > 0 {
            result = append(others[shift:], others[:shift]...)
        } else {
            ss := length - shift - 1
            result = append(others[ss:], others[:ss]...)
        }

        if idx == 0 {
            result = append([]Item{toMix}, result...)
        } else if idx == length-1 {
            result = append(result, toMix)
        } else {
            result = append(result[:idx+1], result[idx:]...)
            result[idx] = toMix
        }
    }

    dist = util.Abs(dist) % length
    //fmt.Printf("%v => shift=%v, move=%v\n", toMix.Value, shift, dist)

    for n := dir * dist; n != 0; n -= dir {
        newIdx := (idx + dir + length) % length
        result[idx], result[newIdx] = result[newIdx], result[idx]
        idx = newIdx
    }

    //fmt.Printf("moved %v@%v to %v -> %v\n", toMix.Value, idx, newIdx, result)
    return result
}

func parseInput(input []string) ([]Item, error) {
    items := []Item{}
    for i, line := range input {
        val, err := strconv.Atoi(line)
        if err != nil {
            return nil, err
        }

        items = append(items, Item{Value: val, Index: i})
    }

    return items, nil
}
