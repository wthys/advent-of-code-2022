package day13

import (
    "strings"
    "fmt"
    "sort"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/solutions/day13/item"
    "github.com/wthys/advent-of-code-2022/solutions/day13/item/parser"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "13"
}

func (s solution) Part1(input []string) (string, error) {
    pairs, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    total := 0
    for i, pair := range pairs {
        if pair.inOrder() {
            total += i + 1
        }
    }

    return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
    pairs, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    items := []item.Item{}

    for _, pair := range pairs {
        items = append(items, pair.Left, pair.Right)
    }

    sep2 := *item.ItemList(*item.ItemValue(2))
    sep6 := *item.ItemList(*item.ItemValue(6))

    items = append(items, sep2, sep6)

    sort.Sort(ItemSlice(items))

    decoderKey := 1
    for i, item := range items {
        if item.Cmp(sep2) == 0 || item.Cmp(sep6) == 0 {
            decoderKey *= i + 1
        }
    }

    return solver.Solved(decoderKey)
}

type (
    Pair struct {
        Left item.Item
        Right item.Item
    }

    ItemSlice []item.Item
)

func (s ItemSlice) Swap(i, j int) {
    s[j], s[i] = s[i], s[j]
}

func (s ItemSlice) Len() int {
    return len(s)
}

func (s ItemSlice) Less(i, j int) bool {
    return s[i].Cmp(s[j]) < 0
}

func (p Pair) inOrder() bool {
    return p.Left.Cmp(p.Right) <= 0
}

func parseInput(input []string) ([]Pair, error) {

    pairs := []Pair{}

    var (
        left *item.Item = nil
    )

    for n, line := range input {
        if len(strings.TrimSpace(line)) == 0 {
            continue
        }

        item, err := parser.NewParser(strings.NewReader(line)).Parse()
        if err != nil {
            return nil, fmt.Errorf("error on line %v: %v", n+1, err)
        }

        if left == nil {
            //fmt.Printf("left: %v\n", item)
            left = item
        } else {
            //fmt.Printf("right: %v\n", item)
            pairs = append(pairs, Pair{*left, *item})
            left = nil
        }
    }

    return pairs, nil
}
