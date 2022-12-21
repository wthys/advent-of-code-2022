package day3

import (
    "fmt"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "3"
}

type Rucksack struct {
    Left string
    Right string
}

func (r Rucksack) Common() string {

    for _, left := range r.Left {
        for _, right := range r.Right {
            if left == right {
                return string(left)
            }
        }
    }

    return ""
}

func (r Rucksack) Contents() string {
    return r.Left + r.Right
}

const priorities string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getPriority(t string) int {
    return strings.Index(priorities, t)+1
}

func makeRucksack(line string) (Rucksack, error) {
    if len(line) % 2 == 1{
        return Rucksack{}, fmt.Errorf("wrong amount of items in rucksack %q", line)
    }

    var (
        left string = line[:len(line)/2]
        right string = line[len(line)/2:]
    )
    return Rucksack{Left: left, Right: right}, nil
}

func parseInput(input []string) ([]Rucksack, error) {
    rucksacks := []Rucksack{}
    for n, line := range input {
        rucksack, err := makeRucksack(line)
        if err != nil {
            return nil, fmt.Errorf("Rucksack on %v: %v", n+1, err)
        }
        rucksacks = append(rucksacks, rucksack)
    }
    return rucksacks, nil
}

func (s solution) Part1(input []string) (string, error) {

    rucksacks, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    total := 0

    for _, rucksack := range rucksacks {
        com := rucksack.Common()
        prio := getPriority(com)
        total += prio
    }

    return solver.Solved(total)
}


func findBadge(r1, r2, r3 Rucksack) string {
    for _, t := range r1.Contents() {
        i2 := strings.Index(r2.Contents(), string(t))
        i3 := strings.Index(r3.Contents(), string(t))
        if i2 >= 0 && i3 >= 0 {
            return string(t)
        }
    }
    return ""
}

func (s solution) Part2(input []string) (string, error) {
    rucksacks, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    total := 0
    for i := 0; i < len(rucksacks); i += 3 {
        r1 := rucksacks[i]
        r2 := rucksacks[i+1]
        r3 := rucksacks[i+2]
        com := findBadge(r1, r2, r3)
        total += getPriority(com)
    }
    return solver.Solved(total)
}
