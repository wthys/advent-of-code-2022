package day8

import (
    "strconv"
    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/location"
    "github.com/wthys/advent-of-code-2022/grid"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "8"
}

func (s solution) Part1(input []string) (string, error) {
    g, err := parseInput(input)

    if err != nil {
        return "", err
    }

    visible := 0

    vis := func(loc location.Location, height int) {
        if isVisible(g, loc, height, location.New(0,-1)) {
            visible += 1
            return
        }
        if isVisible(g, loc, height, location.New(0,1)) {
            visible += 1
            return
        }
        if isVisible(g, loc, height, location.New(-1,0)) {
            visible += 1
            return
        }
        if isVisible(g, loc, height, location.New(1,0)) {
            visible += 1
            return
        }
    }

    g.Apply(vis)

    return strconv.Itoa(visible), nil
}

func (s solution) Part2(input []string) (string, error) {
    g, err := parseInput(input)

    if err != nil {
        return "", err
    }

    best := 0

    scenicScore := func(loc location.Location, height int) {
        up    := viewingDistance(g, loc, height, location.New(0,-1))
        if up == 0 {
            return
        }

        down  := viewingDistance(g, loc, height, location.New(0,1))
        if down == 0 {
            return
        }

        left  := viewingDistance(g, loc, height, location.New(-1,0))
        if left == 0 {
            return
        }

        right := viewingDistance(g, loc, height, location.New(1,0))

        score := up * down * left * right

        if score > best {
            best = score
        }
    }

    g.Apply(scenicScore)

    return strconv.Itoa(best), nil
}


func parseInput(input []string) (*grid.Grid[int], error) {
    g := grid.New[int]()

    for y, line := range input {
        for x, val := range line {
            hgt, err := strconv.Atoi(string(val))
            if err != nil {
                return nil, err
            }
            g.Set(location.New(x, y), hgt)
        }
    }

    return g, nil
}

func isVisible(g *grid.Grid[int], loc location.Location, height int, dir location.Location) bool {
    i := 1
    h, err := g.Get(loc.Add(dir.Scale(i)))
    for err == nil {
        if h >= height {
            return false
        }
        i += 1
        h, err = g.Get(loc.Add(dir.Scale(i)))
    }
    return true
}

func viewingDistance(g *grid.Grid[int], loc location.Location, height int, dir location.Location) int {
    i := 1
    h, err := g.Get(loc.Add(dir.Scale(i)))
    for err == nil {
        if h >= height {
            return i
        }
        i += 1
        h, err = g.Get(loc.Add(dir.Scale(i)))
    }

    if err != nil && i == 1 {
        return 0
    }
    
    return i-1
}
