package day25

import (
    "fmt"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "25"
}

func (s solution) Part1(input []string) (string, error) {
    numbers, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    total := 0
    for _, num := range numbers {
        total += num
    }

    return solver.Solved(IntToSnafu(total))
}

func (s solution) Part2(input []string) (string, error) {
    return solver.NotImplemented()
}

var (
    SNAFUDIGITS = map[rune]int{
        '2': 2,
        '1': 1,
        '0': 0,
        '-': -1,
        '=': -2,
    }
    SNAFUREVERSE = map[int]string {
        2: "2",
        1: "1",
        0: "0",
        -1: "-",
        -2: "=",
    }
)

func SnafuToInt(snafu string) (int, error) {
    num := 0

    for _, ch := range snafu {
        digit, ok := SNAFUDIGITS[ch]
        if !ok {
            return 0, fmt.Errorf("%v is not valid SNAFU digit", string(ch))
        }
        num = 5 * num + digit
    }

    return num, nil
}

func IntToSnafu(num int) string {
    snafu := ""

    rem := 0

    for {
        num, rem = num / 5, num % 5
        if rem <= 2 {
            snafu = SNAFUREVERSE[rem] + snafu
        } else {
            num += 1
            snafu = SNAFUREVERSE[rem-5] + snafu
        }

        if num == 0 {
            break
        }
    }

    return snafu
}


func parseInput(input []string) ([]int, error) {
    numbers := []int{}

    for nr, line := range input {
        num, err := SnafuToInt(line)
        if err != nil {
            return nil, fmt.Errorf("could not parse line %v %q: %v", nr, line, err)
        }

        //fmt.Printf("%6v -> %6v -> %6v\n", line, num, IntToSnafu(num))

        numbers = append(numbers, num)
    }

    return numbers, nil
}
