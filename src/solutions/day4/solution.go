package day4

import (
    "fmt"
    "strings"
    "strconv"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "4"
}


type Range struct {
    lo int
    hi int
}

func (r Range) Contains(o Range) bool {
    return r.lo <= o.lo && r.hi >= o.hi
}

func (r Range) Overlaps(o Range) bool {
    return (r.lo <= o.hi && r.hi >= o.lo)
}

func (r Range) String() string {
    return fmt.Sprintf("%d-%d", r.lo, r.hi)
}

func MakeRange(lo, hi int) Range {

    if lo > hi {
        return Range{hi, lo}
    }

    return Range{lo, hi}
}

func parseInput(input []string) [][]Range {
    pairs := [][]Range{}

    for _, line := range input {
        elves := strings.Split(line, ",")

        elf1 := strings.Split(elves[0], "-")
        elf2 := strings.Split(elves[1], "-")

        lo1, _ := strconv.Atoi(elf1[0])
        hi1, _ := strconv.Atoi(elf1[1])
        lo2, _ := strconv.Atoi(elf2[0])
        hi2, _ := strconv.Atoi(elf2[1])

        pairs = append(pairs, []Range{MakeRange(lo1, hi1), MakeRange(lo2,hi2)})
    }
    return pairs
}


func (s solution) Part1(input []string) (string, error) {
    pairs := parseInput(input)

    total := 0
    for _, pair := range pairs {
        if pair[0].Contains(pair[1]) || pair[1].Contains(pair[0]) {
            total += 1
        }
    }
    return strconv.Itoa(total), nil
}

func (s solution) Part2(input []string) (string, error) {
    pairs := parseInput(input)

    total := 0
    for _, pair := range pairs {
        if pair[0].Overlaps(pair[1]) {
            total += 1
        }
    }
    return strconv.Itoa(total), nil
}
