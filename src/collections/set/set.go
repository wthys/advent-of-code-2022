package set

import (
    "fmt"
    "strings"
)

type (
    empty struct{}

    Set[T comparable] struct {
        contents map[T]empty
    }

    SetDoFunction[T comparable] func(value T) bool
)

func New[T comparable](values ...T) *Set[T] {
    set := Set[T]{map[T]empty{}}
    for _, value := range values {
        set.Add(value)
    }
    return &set
}

func (s Set[T]) String() string {
    str := strings.Builder{}
    fmt.Fprint(&str, "<")
    s.Do(func(value T) bool {
        fmt.Fprintf(&str, " %v", value)
        return true
    })
    fmt.Fprint(&str, " >")
    return str.String()
}

func (s Set[T]) Values() []T {
    vals := []T{}
    s.Do(func(value T) bool {
        vals = append(vals, value)
        return true
    })
    return vals
}

func (s *Set[T]) Add(value T) *Set[T] {
    (*s).contents[value] = empty{}
    return s
}

func (s Set[T]) Has(value T) bool {
    _, ok := s.contents[value]
    return ok
}

func (s Set[T]) Len() int {
    return len(s.contents)
}

func (s *Set[T]) Remove(value T) *Set[T] {
    delete((*s).contents, value)
    return s
}

func (s Set[T]) Intersect(other *Set[T]) *Set[T] {
    common := New[T]()

    s.Do(func(value T) bool {
        if other.Has(value) {
            common.Add(value)
        }
        return true
    })

    return common
}

func (s Set[T]) Union(other *Set[T]) *Set[T] {
    union := New[T]()
    adder := func(value T) bool {
        union.Add(value)
        return true
    }

    s.Do(adder)
    other.Do(adder)

    return union
}

func (s Set[T]) Subtract(other *Set[T]) *Set[T] {
    sub := New[T]()

    s.Do(func(value T) bool {
        if !other.Has(value) {
            sub.Add(value)
        }
        return true
    })

    return sub
}

func (s Set[T]) Do(doer SetDoFunction[T]) {
    for value, _ := range s.contents {
        if !doer(value) {
            break
        }
    }
}
