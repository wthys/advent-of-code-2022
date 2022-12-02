package day2

import (
    "strings"
    "strconv"
    "fmt"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "2"
}

type Round struct {
    Elf string
    You string
}

var (
    scoring = map[string]int{
        "X": 1,
        "Y": 2,
        "Z": 3,
    }
)

const (
    elfMoves = "ABC"
    youMoves = "XYZ"
)

func (r Round) Score() int {
    moveScore := scoring[r.You]

    switch r.Elf {
    case "A": // Rock
        switch r.You {
        case "X": // Rock
            return moveScore+3
        case "Y": // Paper
            return moveScore+6
        case "Z": // Scissors
            return moveScore
        }
    case "B": // Paper
        switch r.You {
        case "X": // Rock
            return moveScore
        case "Y": // Paper
            return moveScore+3
        case "Z": // Scissors
            return moveScore+6
        }
    case "C": // Scissors
        switch r.You {
        case "X": // Rock
            return moveScore+6
        case "Y": // Paper
            return moveScore
        case "Z": // Scissors
            return moveScore+3
        }
    }

    return 0
}

func (r Round) Score2() int {

    switch r.Elf {
    case "A": // Rock
        switch r.You {
        case "X": // Lose
            return 0 + 3
        case "Y": // Draw
            return 3 + 1
        case "Z": // Win
            return 6 + 2
        }
    case "B": // Paper
        switch r.You {
        case "X": // Lose
            return 0 + 1
        case "Y": // Draw
            return 3 + 2
        case "Z": // Win
            return 6 + 3
        }
    case "C": // Scissors
        switch r.You {
        case "X": // Lose
            return 0 + 2
        case "Y": // Draw
            return 3 + 3
        case "Z": // Win
            return 6 + 1
        }
    }

    return 0

}

func MakeRound(elf, you string) (Round, error) {
    if !strings.Contains(elfMoves, elf) {
        return Round{}, fmt.Errorf("Elf cannot make %q move", elf)
    }
    if !strings.Contains(youMoves, you) {
        return Round{}, fmt.Errorf("You cannot make %q move", you)
    }

    return Round{elf, you}, nil
}

func parseRounds(input []string) ([]Round, error) {
    rounds := make([]Round, 0)

    for n, line := range input {
        parts := strings.Split(line, " ")

        if len(parts) < 2 {
            return nil, fmt.Errorf("Invalid round %q on line %v", line, n)
        }

        round, err := MakeRound(parts[0], parts[1])
        if err != nil {
            return nil, err
        }

        rounds = append(rounds, round)
    }

    return rounds, nil
}

func (s solution) Part1(input []string) (string, error) {
    score := 0

    rounds, err := parseRounds(input)
    if err != nil {
        return solver.Unsolved, err
    }

    for _, round := range rounds {
        score += round.Score()
    }

    return strconv.Itoa(score), nil
}

func (s solution) Part2(input []string) (string, error) {
    score := 0

    rounds, err := parseRounds(input)
    if err != nil {
        return solver.Unsolved, err
    }

    for _, round := range rounds {
        score += round.Score2()
    }

    return strconv.Itoa(score), nil
}
