package pathfinding

import (
    "fmt"

    "github.com/wthys/advent-of-code-2022/collections/set"
)

type (
    DistMap[T comparable] map[T]int
    PrevMap[T comparable] map[T]*T

    Dijkstra[T comparable] interface {
        ShortestPathTo(end T) []T
        ShortestPathLengthTo(end T) int
        DoNodes(doer func(node T) bool)
    }

    SimpleDijkstra[T comparable] struct {
        Start T
        dist DistMap[T]
        prev PrevMap[T]
    }

    NeejberFunc[T comparable] func(node T) []T
    ExitFunc[T comparable] func(node T) bool
)

const (
    INFINITE = int((^uint(0)) >> 1)
)

func ControlledDijkstra[T comparable](start T, neejbers NeejberFunc[T], exitters ...ExitFunc[T]) Dijkstra[T] {
    dist := DistMap[T]{}
    prev := PrevMap[T]{}
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

        stop := false
        for _, exit := range exitters {
            if exit(node) {
                stop = true
            }
        }
        if stop { break }

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

    return SimpleDijkstra[T]{start, dist, prev}

}

func ConstructDijkstra[T comparable](start T, neejbers NeejberFunc[T]) Dijkstra[T] {
    dist := DistMap[T]{}
    prev := PrevMap[T]{}
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

    return SimpleDijkstra[T]{start, dist, prev}
}

func (d SimpleDijkstra[T]) ShortestPathTo(end T) []T {
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

func (d SimpleDijkstra[T]) DoNodes(doer func(node T) bool) {
    for node, _ := range d.prev {
        if !doer(node) {
            return
        }
    }
}


func (d SimpleDijkstra[T]) ShortestPathLengthTo(end T) int {
    dist, ok := d.dist[end]
    if !ok {
        return INFINITE
    }
    return dist
}

func ShortestPath[T comparable](start, end T, neejbers NeejberFunc[T]) ([]T, error) {

    d := ControlledDijkstra(start, neejbers, func (node T) bool { return node == end })

    path := d.ShortestPathTo(end)
    if path == nil {
        return nil, fmt.Errorf("could not find a path from %v to %v", start, end)
    }

    return path, nil
}

func closest[T comparable](Q *set.Set[T], dist DistMap[T]) (T, error) {
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
