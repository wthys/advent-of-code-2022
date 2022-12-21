package pathfinding

import (
    "fmt"

    "github.com/wthys/advent-of-code-2022/collections/set"
)

type (
    distMap[T comparable] map[T]int
    prevMap[T comparable] map[T]*T

    Dijkstra[T comparable] struct {
        Start T
        dist distMap[T]
        prev prevMap[T]
    }

    NeejberFunc[T comparable] func(node T) []T
)

const (
    INFINITE = int((^uint(0)) >> 1)
)

func ConstructDijkstra[T comparable](nodes []T, start T, neejbers NeejberFunc[T]) Dijkstra[T] {
    dist := distMap[T]{}
    prev := prevMap[T]{}
    queue := []T{}
    visited := set.New[T]()

    for _, loc := range nodes {
        dist[loc] = INFINITE
        prev[loc] = nil
        queue = append(queue, loc)
    }

    dist[start] = 0

    for len(queue) > 0 {
        i, node := closest(queue, dist)
        queue = append(queue[:i], queue[i+1:]...)
        visited.Add(node)

        for _, neejber := range neejbers(node) {
            if visited.Has(neejber) {
                continue
            }
            alt := dist[node] + 1
            if alt < dist[neejber] {
                dist[neejber] = alt
                prev[neejber] = &node
            }
        }
    }

    return Dijkstra[T]{start, dist, prev}
}

func (d Dijkstra[T]) ShortestPathTo(end T) []T {
    path := []T{}
    node := &end
    for *node != d.Start && node != nil {
        path = append([]T{*node}, path...)
        node = d.prev[*node]
    }

    if node == nil {
        return nil
    }

    return path
}

func (d Dijkstra[T]) ShortestPathLengthTo(end T) int {
    return d.dist[end]
}

func ShortestPath[T comparable](nodes []T, start, end T, neejbers NeejberFunc[T]) ([]T, error) {

    d := ConstructDijkstra(nodes, start, neejbers)

    path := d.ShortestPathTo(end)
    if path == nil {
        return nil, fmt.Errorf("could not find a path from %v to %v", start, end)
    }

    return path, nil
}

func closest[T comparable](Q []T, dist distMap[T]) (int, T) {
    shortest := INFINITE
    si := -1
    sloc := *new(T)

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
