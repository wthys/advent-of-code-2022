package day14

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/grid"
    "github.com/wthys/advent-of-code-2022/util"
    "github.com/wthys/advent-of-code-2022/location"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "14"
}

func (s solution) Part1(input []string) (string, error) {

    paths, err := parseInput(input)
    if err != nil {
        return "", err
    }

    g := prepareGrid(paths)

    source := location.New(500,0)

    g.Set(source, "+")
    //fmt.Println("== initial state ==")
    //g.Print()

    n := 0
    for {
        rest, err := dropSand(g, source)
        if err != nil {
            break
        }

        n += 1
        g.Set(rest, "o")

        //fmt.Printf("== sand #%v ==\n", n)
        //g.Print()
    }

    return strconv.Itoa(n), nil
}

func (s solution) Part2(input []string) (string, error) {

    paths, err := parseInput(input)
    if err != nil {
        return "", err
    }

    g := prepareGrid(paths)

    source := location.New(500,0)
    g.Set(source, "+")

    //floor
    bounds, _ := g.Bounds()
    for x := bounds.Xmin - bounds.Ymax - 2; x <= bounds.Xmax + bounds.Ymax + 2; x++ {
        loc := location.New(x, bounds.Ymax + 2)
        g.Set(loc, "=")
    }

    n := 0
    for {
        rest, err := dropSand(g, source)
        if rest == source || err != nil {
            break
        }

        n += 1
        g.Set(rest, "o")
    }

    return strconv.Itoa(n+1), nil
}

type (
    Path []location.Location
)

var (
    errAbyss = fmt.Errorf("sand fell into the abyss")
)

func dropSand(g *grid.Grid[string], source location.Location) (location.Location, error) {
    bounds, err := g.Bounds()
    if err != nil {
        return location.Location{}, errAbyss
    }

    dirs := []location.Location{
        location.New(0,1),
        location.New(-1,1),
        location.New(1,1),
    }

    loc := source
    for bounds.Contains(loc) {
        newLoc := loc
        for _, dir := range dirs {
            cand := loc.Add(dir)
            val, _ := g.Get(cand)
            if val == "." {
                newLoc = cand
                break
            }
        }

        if loc == newLoc {
            return newLoc, nil
        }

        loc = newLoc
    }

    return location.Location{}, errAbyss
}

func prepareGrid(paths []Path) *grid.Grid[string] {
    g := grid.WithDefault(".")

    for _, path := range paths {
        for i, pos := range path {
            if (i+1 == len(path)) {
                continue
            }
            diff := path[i+1].Subtract(pos)
            length := util.Abs(diff.X + diff.Y)
            for k := 0; k <= length; k++ {
                loc := pos.Add(diff.Unit().Scale(k))
                g.Set(loc, "#")
            }
        }
    }

    return g
}

func parseInput(input []string) ([]Path, error) {
    reLoc := regexp.MustCompile("(\\d+),(\\d+)")

    paths := []Path{}

    for _, line := range input {
        locs := reLoc.FindAllStringSubmatch(line, -1)
        if locs == nil {
            return nil, fmt.Errorf("could not parse %q", line)
        }

        path := Path{}
        for _, loc := range locs {
            x, _ := strconv.Atoi(loc[1])
            y, _ := strconv.Atoi(loc[2])

            path = append(path, location.New(x, y))
        }
        paths = append(paths, path)
    }

    return paths, nil
}

func (p Path) String() string {
    str := strings.Builder{}
    for i, loc := range p {
        if i > 0 {
            fmt.Fprint(&str, " -> ")
        }
        fmt.Fprint(&str, loc.String())
    }
    return str.String()
}
