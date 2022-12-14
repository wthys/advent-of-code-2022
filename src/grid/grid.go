package grid

import (
    "fmt"
    "github.com/wthys/advent-of-code-2022/location"
)

type (
    Grid[T any] struct {
        defaultFunc DefaultFunction[T]
        data map[location.Location]T
    }

    Bounds struct {
        Xmin, Xmax, Ymin, Ymax int
    }

    DefaultFunction[T any] func(loc location.Location) (T, error)
    ApplyFunction[T any] func(loc location.Location, value T)
)

// `DefaultValue` creates a `DefaultFunction` that always returns the provided
// value.
func DefaultValue[T any](value T) DefaultFunction[T] {
    return func(_ location.Location) (T, error) {
        return value, nil
    }
}

// `DefaultZero` creates a `DefaultFunction` that always returns the 'zero' value.
func DefaultZero[T any]() DefaultFunction[T] {
    return DefaultValue(*new(T))
}

// `DefaultError` creates a `DefaultFunction` that always returns an error "no
// value at <loc>"
func DefaultError[T any]() DefaultFunction[T] {
    return func(loc location.Location) (T, error) {
        return *new(T), fmt.Errorf("no value at %v", loc)
    }
}

// `New` creates a `Grid` using the `DefaultError` `DefaultFunction` for unknown
// `Location`s. Equivalent to `WithDefaultFunc(DefaultError())`.
func New[T any]() *Grid[T] {
    return WithDefaultFunc(DefaultError[T]())
}

// `WithDefault` creates a `Grid` using the `DefaultValue` `DefaultFunction` for
// unknown `Location`s. Equivalent to `WithDefaultFunc(DefaultValue(value))`.
func WithDefault[T any](value T) *Grid[T] {
    return WithDefaultFunc(DefaultValue(value))
}

// `WithDefaultFunc` creates a `Grid` using the provided `DefaultFunction` for
// unknown `Location`s.
func WithDefaultFunc[T any](defaultFunc DefaultFunction[T]) *Grid[T] {
    return &Grid[T]{defaultFunc, map[location.Location]T{}}
}

// `Get` retrieves the value stored at `loc`. If there is no value stored, the
// `Grid`'s `DefaultFunction` is called. If no `DefaultFunction` was set,
// `DefaultError[T]()` is used.
func (g *Grid[T]) Get(loc location.Location) (T, error) {
    val, ok := g.data[loc]
    if ok {
        return val, nil
    }
    if g.defaultFunc == nil {
        return DefaultError[T]()(loc)
    }

    return g.defaultFunc(loc)
}

// `Set` stores a value at `loc`.
func (g *Grid[T]) Set(loc location.Location, value T) {
    g.data[loc] = value
}

// `Remove` removes the stored value at `loc`, if any.
func (g *Grid[T]) Remove(loc location.Location) {
    delete(g.data, loc)
}

// `Apply` applies a function to all stored values. Both the `Location` and the
// value are provided to the given `ApplyFunction`.
func (g *Grid[T]) Apply(applyFunc ApplyFunction[T]) {
    for loc, value := range g.data {
        applyFunc(loc, value)
    }
}

// `Bounds` finds the bounding box of the `Location`s of the stored values.
// Returns an error when there are no stored values.
func (g *Grid[T]) Bounds() (Bounds, error) {
    if len(g.data) == 0 {
        return Bounds{}, fmt.Errorf("no values in grid")
    }

    bounds := Bounds{0,0,0,0}
    found := false
    apply := func(loc location.Location, _ T) {
        if !found {
            bounds.Xmin = loc.X
            bounds.Xmax = loc.X
            bounds.Ymin = loc.Y
            bounds.Ymax = loc.Y
            found = true
            return
        }

        if loc.X < bounds.Xmin {
            bounds.Xmin = loc.X
        }
        if loc.X > bounds.Xmax {
            bounds.Xmax = loc.X
        }

        if loc.Y < bounds.Ymin {
            bounds.Ymin = loc.Y
        }
        if loc.Y > bounds.Ymax {
            bounds.Ymax = loc.Y
        }
    }
    g.Apply(apply)

    return bounds, nil
}

// `Len` returns the number of stored values.
func (g *Grid[T]) Len() int {
    return len(g.data)
}

func (g *Grid[T]) Print() {
    bounds, err := g.Bounds()

    if err != nil {
        fmt.Println()
        return
    }

    for y := bounds.Ymin; y <= bounds.Ymax; y++ {
        for x := bounds.Xmin; x <= bounds.Xmax; x++ {
            pos := location.New(x, y)
            val, err := g.Get(pos)
            if err != nil {
                fmt.Print(".")
                continue
            }
            fmt.Print(val)
        }
        fmt.Println()
    }
}

func (b *Bounds) Contains(loc location.Location) bool {
    return loc.X >= b.Xmin && loc.X <= b.Xmax && loc.Y >= b.Ymin && loc.Y <= b.Ymax
}

func (b Bounds) Width() int {
    return b.Xmax - b.Xmin + 1
}

func (b Bounds) Height() int {
    return b.Ymax - b.Ymin + 1
}

func (b Bounds) Apply(apply func(loc location.Location)) {
    for y := b.Ymin; y <= b.Ymax; y++ {
        for x := b.Xmin; x <= b.Xmax; x++ {
            loc := location.New(x, y)
            apply(loc)
        }
    }
}
