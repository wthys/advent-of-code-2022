package day16

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/golang-collections/collections/set"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/util"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "16"
}

func (s solution) Part1(input []string) (string, error) {
    tunnels, rates, err := parseInput(input)
    if err != nil {
        return "", err
    }

    paths := calculatePaths(tunnels, rates)

    pressure, path := findBestPressurePath("AA", paths, rates, 30, set.New())

    return strconv.Itoa(pressure), nil
}

func (s solution) Part2(input []string) (string, error) {
    return solver.NotImplemented()
}

func findBestPressurePath(start string, paths shortestPaths, rates ValveRates, timeLeft int, visited *set.Set) (pressure int, path []string) {
    if timeLeft <= 0 {
        return 0, []string{}
    }

    if visited.Has(start) {
        return 0, []string{}
    }

    rate := rates[start]
    if rate <= 0 && start != "AA" {
        return 0, []string{}
    }
    newVisited := set.New(start).Union(opened)

    valveOpened := util.IIf(rate > 0, 1, 0)

    bestPressure := 0
    bestPath := []string{}
    for next, edge := range paths[start] {
        if visited.Has(next) {
            continue
        }

        time := len(edge)
        pres, pth := findBestPressurePath(next, paths, rates, timeLeft - time - valveOpened, newVisited)
        if pres > bestPressure {
            bestPressure = pres
            bestPath = pth
        }
    }

    pressure = bestPressure + (timeLeft - valveOpened) * rate
    path = append([]string{start}, bestPath...)


    return
}



type (
    TunnelMap map[string][]string
    ValveRates map[string]int
)

func calculatePaths(tunnels TunnelMap, rates ValveRates) shortestPaths {
    relevantValves := []string{}
    allValves := []string{}
    for valve, rate := range rates {
        allValves = append(allValves, valve)
        if rate > 0 {
            relevantValves = append(relevantValves, valve)
        }
    }

    neejbers := func(node string) []string {
        return tunnels[node]
    }

    paths := shortestPaths{}
    for _, start := range append([]string{"AA"}, relevantValves...) {
        _, ok := paths[start]
        if !ok {
            paths[start] = map[string][]string{}
        }

        _, prev := dijkstra(allValves, start, neejbers)
        for _, end := range relevantValves {
            if start == end {
                continue
            }
            path := shortestPathFromDijkstra(prev, start, end)
            if path == nil {
                continue
            }

            paths[start][end] = path
        }

    }

    return paths
}

func parseInput(input []string) (TunnelMap, ValveRates, error) {
    tunnels := TunnelMap{}
    rates := ValveRates{}

    reRate := regexp.MustCompile("-?\\d+")
    reValves := regexp.MustCompile("[A-Z]{2}")

    for nr, line := range input {
        valves := reValves.FindAllString(line, -1)
        if valves == nil {
            return nil, nil, fmt.Errorf("could not find valve names on line %v: %q", nr, line)
        }

        tunnels[valves[0]] = valves[1:]

        rate := reRate.FindString(line)
        if rate == "" {
            return nil, nil, fmt.Errorf("could not find valve rate on line %v: %q", nr, line)
        }

        rateValue, _ := strconv.Atoi(rate)
        rates[valves[0]] = rateValue
    }

    return tunnels, rates, nil
}

type (
    distMap map[string]int
    prevMap map[string]string
    shortestPaths map[string](map[string][]string)
)

var (
    undefined = ""
    infinite = 1_000_000_000
)

func dijkstra(nodes []string, start string, neejbers func(string) []string) (distMap, prevMap) {
    dist := distMap{}
    prev := prevMap{}
    queue := []string{}
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

func shortestPathFromDijkstra(prev prevMap, start, end string) []string {
    path := []string{}
    node := end
    for node != start && node != undefined {
        path = append([]string{node}, path...)
        node = prev[node]
    }
    if node == undefined {
        return nil
    }

    return path
}

func findShortestPath(nodes []string, start, end string, neejbers func(string) []string) ([]string, error) {

    _, prev := dijkstra(nodes, start, neejbers)

    path := shortestPathFromDijkstra(prev, start, end)
    if path == nil {
        return nil, fmt.Errorf("could not find a path from %v to %v", start, end)
    }

    return path, nil
}

func closest(Q []string, dist distMap) (int, string) {
    shortest := infinite+1
    si := -1
    sloc := ""

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
