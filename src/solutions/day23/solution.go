package day23

import (
    "fmt"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/location"
    "github.com/wthys/advent-of-code-2022/grid"
    "github.com/wthys/advent-of-code-2022/collections/set"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "23"
}

func (s solution) Part1(input []string) (string, error) {
    grove, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    //fmt.Println("== Initial State ==")
    //fmt.Println(grove)

    for round := 1; round <= 10; round++ {
        //fmt.Printf("== End of Round %v ==\n", round)
        done := grove.Round()
        //fmt.Println(grove)
        if done {
            break
        }
    }

    count := grove.CountInRectangle()

    return solver.Solved(count["."])
}

func (s solution) Part2(input []string) (string, error) {
    grove, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    //fmt.Println("== Initial State ==")
    //fmt.Println(grove)

    round := 1
    for {
        //fmt.Printf("== End of Round %v ==\n", round)
        done := grove.Round()
        //fmt.Println(grove)
        if done {
            break
        }
        round += 1
    }

    return solver.Solved(round)
}

const (
    NORTH = "N"
    NORTHEAST = "NE"
    EAST = "E"
    SOUTHEAST = "SE"
    SOUTH = "S"
    SOUTHWEST = "SW"
    WEST = "W"
    NORTHWEST = "NW"
    NOWHERE = "X"
)

var (
    DIRECTIONS = map[string]location.Location{
        NORTH: {0,-1},
        NORTHEAST: {1,-1},
        EAST: {1,0},
        SOUTHEAST: {1,1},
        SOUTH: {0,1},
        SOUTHWEST: {-1,1},
        WEST: {-1,0},
        NORTHWEST: {-1,-1},
    }
)

type (
    Elf struct {
        Pos location.Location
    }

    Grove struct {
        Elves []Elf
        Checkers map[string]DirChecker
        Grid *grid.Grid[string]
        Opts []string
    }

    DirChecker func(location.Location) string
)

func (e *Elf) MoveTo(pos location.Location) {
    (*e).Pos = pos
}

func (g *Grove) Round() bool {
    //phase 1
    proposals := map[location.Location]*set.Set[int]{}

    for idx, elf := range (*g).Elves {
        proposedDir := ""
        for _, dir := range (*g).Opts {
            cand := (*g).Checkers[dir](elf.Pos)
            if cand != "" {
                proposedDir = cand
                break
            }
        }

        if proposedDir == "" || proposedDir == NOWHERE {
            continue
        }

        proposal := elf.Pos.Add(DIRECTIONS[proposedDir])
        _, ok := proposals[proposal]
        if !ok {
            proposals[proposal] = set.New[int]()
        }
        proposals[proposal].Add(idx)
    }

    if len(proposals) == 0 {
        return true
    }

    //phase 2
    for loc, elves := range proposals {
        if elves.Len() == 1 {
            elves.Do(func (idx int) bool {
                (*g).Grid.Remove((*g).Elves[idx].Pos)
                (*g).Grid.Set(loc, "#")
                (*g).Elves[idx].MoveTo(loc)
                return false
            })
        }
    }

    //round end
    opts := (*g).Opts
    (*g).Opts = append(opts[1:], opts[0])

    return false
}

func (g Grove) CountInRectangle() map[string]int {
    count := map[string]int{"#": 0, ".": 0}
    b, err := g.Grid.Bounds()
    if err != nil {
        return count
    }

    count["#"] = len(g.Elves)
    count["."] = (b.Width() * b.Height()) - len(g.Elves)
    return count
}

func (g Grove) String() string {
    b, err := g.Grid.Bounds()
    if err != nil {
        return ""
    }

    str := strings.Builder{}

    for y := b.Ymin; y <= b.Ymax; y++ {
        for x := b.Xmin; x <= b.Xmax; x++ {
            loc := location.New(x, y)
            val, _ := g.Grid.Get(loc)
            fmt.Fprint(&str, val)
        }
        fmt.Fprint(&str, "\n")
    }

    return str.String()
}

func MakeChecker(g *grid.Grid[string], result string, directions ...string) DirChecker {
    return func(loc location.Location) string {
        count := 0
        for _, dir := range DIRECTIONS {
            val, err := g.Get(loc.Add(dir))
            if err == nil && val == "#" {
                count += 1
            }
        }
        if count == 0 {
            return NOWHERE
        }

        for _, dir := range directions {
            val, err := g.Get(loc.Add(DIRECTIONS[dir]))
            if err == nil && val == "#" {
                return ""
            }
        }
        return result
    }
}

func parseInput(input []string) (Grove, error) {
    elves := []Elf{}

    g := grid.WithDefault(".")

    for y, line := range input {
        for x, ch := range line {
            if string(ch) != "#" {
                continue
            }

            pos := location.New(x, y)
            g.Set(pos, "#")
            elves = append(elves, Elf{pos})
        }
    }

    checkers := map[string]DirChecker{
        NORTH: MakeChecker(g, NORTH, NORTHEAST, NORTH, NORTHWEST),
        SOUTH: MakeChecker(g, SOUTH, SOUTHEAST, SOUTH, SOUTHWEST),
        WEST: MakeChecker(g, WEST, NORTHWEST, WEST, SOUTHWEST),
        EAST: MakeChecker(g, EAST, NORTHEAST, EAST, SOUTHEAST),
    }

    return Grove{elves, checkers, g, []string{"N", "S", "W", "E"}}, nil
}
