package day1

import (
    "fmt"
    "strconv"
    "sort"

    "github.com/wthys/advent-of-code-2022/solver"
)


type solution struct {}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "1"
}


type Elf struct {
    Rations []int
}

func (e Elf) TotalCalories() int {
    if e.Rations == nil {
        return 0
    }

    total := 0
    for _, ration := range e.Rations {
        total += ration
    }

    return total
}


func parseInput(lines []string) ([]Elf, error) {
    elves := make([]Elf, 0)

    rations := make([]int, 0)

    for nr, line := range lines {
        if len(line) == 0 {
            elf := Elf{ rations }
            elves = append(elves, elf)
            rations = make([]int, 0)
	} else {
            energy, err := strconv.Atoi(line)
            if err != nil {
                return nil, fmt.Errorf("value on line #%v is not a number (%v)", nr+1, line)
            }

            rations = append(rations, energy)
        }
    }

    // don't forget the last Elf
    elves = append(elves, Elf{ rations })

    return elves, nil
}

func (s solution) Part1(input []string) (string, error) {

    elves, err := parseInput(input)

    if err != nil {
          return solver.Error(err)
    }

    var best int = elves[0].TotalCalories()

    for _, elf := range elves[1:] {

        total := elf.TotalCalories()
        if total > best {
            best = total
        }

    }

    return solver.Solved(best)
}

func (s solution) Part2(input []string) (string, error) {

    elves, err := parseInput(input)

    if err != nil {
        return solver.Error(err)
    }

    top := elves[:]
    sort.SliceStable(top, func(i, j int) bool {
        return top[i].TotalCalories() > top[j].TotalCalories()
    })

    total := top[0].TotalCalories() + top[1].TotalCalories() + top[2].TotalCalories()

    return solver.Solved(total)
}
