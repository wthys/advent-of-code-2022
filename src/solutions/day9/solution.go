package day9

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/wthys/advent-of-code-2022/solver"
    loc "github.com/wthys/advent-of-code-2022/location"
    "github.com/wthys/advent-of-code-2022/grid"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "9"
}

func (s solution) Part1(input []string) (string, error) {
    moves, err := parseInput(input)
    if err != nil {
        return "", err
    }

    knots := []loc.Location{loc.New(0,0), loc.New(0,0)}

    total := Simulate(knots, moves)

    return strconv.Itoa(total), nil
}

func (s solution) Part2(input []string) (string, error) {
    moves, err := parseInput(input)
    if err != nil {
        return "", err
    }

    knots := []loc.Location{}
    for len(knots) < 10 {
        knots = append(knots, loc.New(0,0))
    }

    total := Simulate(knots, moves)

    return strconv.Itoa(total), nil
}

type Move struct {
    Heading string
    Amount int
}

const (
    UP = "U"
    DOWN = "D"
    LEFT = "L"
    RIGHT = "R"
)

var DIRECTIONS = map[string]loc.Location {
    UP: loc.New(0,-1),
    DOWN: loc.New(0,1),
    LEFT: loc.New(-1,0),
    RIGHT: loc.New(1,0),
}

func Simulate(knots []loc.Location, moves []Move) int {
    g := grid.WithDefault[int](0)

    tail := len(knots)-1

    GridInc(g, knots[tail], 1)

    for _, move := range moves {
        dir := DIRECTIONS[move.Heading]
        for i := 0; i < move.Amount; i++ {
            for k, knot := range knots {
                switch k {
                case 0: // head
                    knots[k] = knot.Add(dir)
                    GridInc(g, knots[k], 0)
                    continue
                case tail: // tail
                    prev := knots[k-1]
                    newPos := moveCloser(knot, prev)
                    knots[k] = newPos
                    GridInc(g, knots[k], 1)
                default: // body
                    prev := knots[k-1]
                    newPos := moveCloser(knot, prev)
                    if newPos == knot {
                        break // don't bother with the rest
                    }
                    knots[k] = newPos
                    GridInc(g, knots[k], 0)
                }
            }
        }
    }

    total := 0
    counter := func(_ loc.Location, n int) {
        if n > 0 {
            total += 1
        }
    }
    g.Apply(counter)

    return total
}

func moveCloser(l, tgt loc.Location) loc.Location {
    diff := tgt.Subtract(l)

    xclose := diff.X >= -1 && diff.X <= 1
    yclose := diff.Y >= -1 && diff.Y <= 1

    if xclose && yclose {
        return l
    }

    return l.Add(diff.Unit())
}

func GridInc(g *grid.Grid[int], l loc.Location, amount int) {
    val, err := g.Get(l)
    if err != nil {
        panic(err)
    }

    g.Set(l, val + amount)
}

func PrintGrid(g *grid.Grid[int], head loc.Location, tail loc.Location) {
    bounds, err := g.Bounds()
    if err != nil {
        panic(err)
    }

    for y := bounds.Ymin; y <= bounds.Ymax; y++{
        for x := bounds.Xmin; x <= bounds.Xmax; x++ {
            l := "."
            pos := loc.New(x, y)
            switch pos {
            case head:
                l = "H"
            case tail:
                l = "T"
            default:
                v, _ := g.Get(pos)
                if v > 0 {
                    l = "#"
                }
            }
            fmt.Print(l)
        }
        fmt.Println()
    }
}

func (m Move) String() string {
    return fmt.Sprintf("%v %v", m.Heading, m.Amount)
}


func parseInput(input []string) ([]Move, error) {
    reMove := regexp.MustCompile("([UDLR])\\s+(\\d+)")

    moves := []Move{}

    for _, line := range input {
        caps := reMove.FindStringSubmatch(line)
        if caps == nil {
            return nil, fmt.Errorf("Could make sense of %q", line)
        }

        heading := caps[1]
        amount, err := strconv.Atoi(caps[2])

        if err != nil {
            return nil, err
        }

        moves = append(moves, Move{heading, amount})

    }

    return moves, nil
}
