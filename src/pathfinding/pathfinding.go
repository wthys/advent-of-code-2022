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

func ConstructDijkstra[T comparable](start T, neejbers NeejberFunc[T]) Dijkstra[T] {
    dist := distMap[T]{}
    prev := prevMap[T]{}
    visited := set.New[T]()

    prev[start] = nil
    dist[start] = 0
    queue := set.New(start)

    for queue.Len() > 0 {
        node, err := closest(queue, dist)
        if err != nil {
            panic(err)
        }
        queue = queue.Remove(node)
        visited.Add(node)

        for _, neejber := range neejbers(node) {
            if visited.Has(neejber) {
                continue
            }
            queue.Add(neejber)
            alt := dist[node] + 1
            ndist, ok := dist[neejber]
            if !ok || alt < ndist {
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
    for node != nil && *node != d.Start {
        path = append([]T{*node}, path...)
        ok := true
        node, ok = d.prev[*node]
        if !ok {
            node = nil
        }
    }

    if node == nil {
        return nil
    }

    return path
}

func (d Dijkstra[T]) DoNodes(doer func(node T) bool) {
    for node, _ := range d.prev {
        if !doer(node) {
            return
        }
    }
}


func (d Dijkstra[T]) ShortestPathLengthTo(end T) int {
    dist, ok := d.dist[end]
    if !ok {
        return INFINITE
    }
    return dist
}

func ShortestPath[T comparable](start, end T, neejbers NeejberFunc[T]) ([]T, error) {

    d := ConstructDijkstra(start, neejbers)

    path := d.ShortestPathTo(end)
    if path == nil {
        return nil, fmt.Errorf("could not find a path from %v to %v", start, end)
    }

    return path, nil
}

func closest[T comparable](Q *set.Set[T], dist distMap[T]) (T, error) {
    shortest := INFINITE
    snode := *new(T)
    found := false
    Q.Do(func(node T) bool {
        d := dist[node]
        if d < shortest {
            shortest = d
            snode = node
            found = true
        }
        return true
    })

    if !found {
        return snode, fmt.Errorf("could not find closest node")
    }

    return snode, nil
}
