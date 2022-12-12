package day12

import (
    "fmt"
    "strings"
    "github.com/golang-collections/collections/set"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/grid"
    "github.com/wthys/advent-of-code-2022/location"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "12"
}

func (s solution) Part1(input []string) (string, error) {
    g := parseInput(input)

    start, end := findStartAndEnd(g)

    nodes := []location.Location{}

    g.Apply(func(loc location.Location, _ string) {
        nodes = append(nodes, loc)
    })

    neejbers := func(loc location.Location) []location.Location {
        return findNeejbersUp(g, loc)
    }

    path, err := findShortestPath(nodes, start, end, neejbers)
    if err != nil {
        return "", err
    }

    //fmt.Println("-= part1 =-")
    //printPath(g, path)

    return fmt.Sprint(len(path)), nil
}

func (s solution) Part2(input []string) (string, error) {
    g := parseInput(input)

    start, end := findStartAndEnd(g)

    candidates := []location.Location{start}
    nodes := []location.Location{}
    g.Apply(func(loc location.Location, elevation string) {
        nodes = append(nodes, loc)
        if elevation == "a" {
            candidates = append(candidates, loc)
        }
    })

    neejbers := func(loc location.Location) []location.Location {
        return findNeejbersDown(g, loc)
    }

    dist, prev := dijkstra(nodes, end, neejbers)

    shortest := infinite+1
    sloc := undefined
    for _, cand := range candidates {
        if dist[cand] < shortest {
            shortest = dist[cand]
            sloc = cand
        }
    }


    path := shortestPathFromDijkstra(prev, end, sloc)
    if path == nil {
        return "", fmt.Errorf("no path found")
    }

    rpath := []location.Location{}
    for _, loc := range path {
        rpath = append([]location.Location{loc}, rpath...)
    }

    //fmt.Println("-= part2 =-")
    //printPath(g, rpath)

    return fmt.Sprint(shortest), nil
}

func printPath(g *grid.Grid[string], path []location.Location) {
    bounds, _ := g.Bounds()

    dirs := map[location.Location]string{
        location.New(1,0): ">",
        location.New(0,1): "v",
        location.New(-1,0): "<",
        location.New(0,-1): "^",
    }

    for y := bounds.Ymin; y <= bounds.Ymax; y++ {
        for x := bounds.Xmin; x <= bounds.Xmax; x++ {
            v := "."
            loc := location.New(x,y)
            for i, el := range path {
                if el != loc {
                    continue
                }
                if i+1 >= len(path) {
                    v = "E"
                    break
                }

                diff := path[i+1].Subtract(loc)
                v = fmt.Sprintf("\033[32m%v\033[0m", dirs[diff])
                break
            }

            if v == "." {
                v, _ = g.Get(loc)
            }

            fmt.Print(v)
        }
        fmt.Println()
    }
}

type (
    distMap map[location.Location]int
    prevMap map[location.Location]location.Location
)

var (
    undefined = location.New(-1337,-1337)
    infinite = 1_000_000_000
)

func dijkstra(nodes []location.Location, start location.Location, neejbers func(location.Location) []location.Location) (distMap, prevMap) {
    dist := distMap{}
    prev := prevMap{}
    queue := []location.Location{}
    visited := set.New()

    for _, loc := range nodes {
        dist[loc] = infinite
        prev[loc] = undefined
        queue = append(queue, loc)
    }

    dist[start] = 0

    for len(queue) > 0 {
        i, node := closest(queue, dist)
        queue = append(queue[:i], queue[i+1:]...)
        visited.Insert(node)

        for _, neejber := range neejbers(node) {
            if visited.Has(neejber) {
                continue
            }
            alt := dist[node] + 1
            if alt < dist[neejber] {
                dist[neejber] = alt
                prev[neejber] = node
            }
        }
    }

    return dist, prev
}

func shortestPathFromDijkstra(prev prevMap, start, end location.Location) []location.Location {
    path := []location.Location{}
    node := end
    for node != start && node != undefined {
        path = append([]location.Location{node}, path...)
        node = prev[node]
    }
    if node == undefined {
        return nil
    }

    return path
}

func findShortestPath(nodes []location.Location, start, end location.Location, neejbers func(location.Location) []location.Location) ([]location.Location, error) {

    _, prev := dijkstra(nodes, start, neejbers)

    path := shortestPathFromDijkstra(prev, start, end)
    if path == nil {
        return nil, fmt.Errorf("could not find a path from %v to %v", start, end)
    }

    return path, nil
}

func closest(Q []location.Location, dist map[location.Location]int) (int, location.Location) {
    shortest := infinite+1
    si := -1
    sloc := location.New(0,0)

    for i, loc := range Q {
        d := dist[loc]
        if d < shortest {
            shortest = d
            si = i
            sloc = loc
        }
    }

    return si, sloc
}

func heightDiff(start, target string) int {
    alfa := "SabcdefghijklmnopqrstuvwxyzE"

    s := strings.Index(alfa, start)
    t := strings.Index(alfa, target)
    switch {
    case s < 0 || t < 0:
        return 1_000_000_000
    default:
        return t - s
    }
}

func findNeejbersUp(g *grid.Grid[string], loc location.Location) []location.Location {
    valid := func(src, tgt string) bool {
        return heightDiff(src, tgt) <= 1
    }
    return identifyNeejbers(g, loc, valid)
}

func findNeejbersDown(g *grid.Grid[string], loc location.Location) []location.Location {
    valid := func(src, tgt string) bool {
        return heightDiff(src, tgt) >= -1
    }
    return identifyNeejbers(g, loc, valid)
}

func identifyNeejbers(g *grid.Grid[string], loc location.Location, comp func(a, b string) bool) []location.Location {
    neejbers := []location.Location{ { 0,-1}, {-1, 0}, { 1, 0}, { 0, 1}, }

    valid := []location.Location{}

    height, err := g.Get(loc)
    if err != nil {
        return nil
    }

    for _, off := range neejbers {
        neejber := loc.Add(off)

        val, err := g.Get(neejber)
        if err != nil {
            continue
        }

        if comp(height, val) {
            valid = append(valid, neejber)
        }
    }
    return valid
}

func findStartAndEnd(g *grid.Grid[string]) (location.Location, location.Location) {
    start := location.New(0,0)
    end := location.New(0,0)

    search := func(loc location.Location, value string) {
        switch value {
            case "S":
                start = loc
            case "E":
                end = loc
            default:
                // nothing
        }
    }

    g.Apply(search)

    return start, end
}

func parseInput(input []string) *grid.Grid[string] {
    g := grid.WithDefaultFunc(grid.DefaultError[string]())

    for y, line := range input {
        for x, ch := range line {
            g.Set(location.New(x, y), string(ch))
        }
    }

    return g
}
