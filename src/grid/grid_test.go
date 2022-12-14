package grid

import (
    "fmt"
    "testing"
    "github.com/wthys/advent-of-code-2022/location"
)

func TestDefaultZeroInt(t *testing.T) {
    def := DefaultZero[int]()
    loc := location.New(1,3)

    want := 0

    val, err := def(loc)
    if err != nil || val != want {
        t.Fatalf("def(%v) = %v, %v, want %v, %v", loc, val, err, want, nil)
    }
}

func TestDefaultZeroFloat(t *testing.T) {
    def := DefaultZero[float64]()
    loc := location.New(4,-1)

    want := 0.0

    val, err := def(loc)
    if err != nil || val != want {
        t.Fatalf("def(%v) = %v, %v, want %v, %v", loc, val, err, want, nil)
    }
}

func TestDefaultValueInt(t *testing.T) {
    def := DefaultValue(42)
    loc := location.New(1, -24)

    want := 42

    val, err := def(loc)
    if err != nil || val != want {
        t.Fatalf("def(%v) = %v, %v, want %v, %v", loc, val, err, want, nil)
    }
}

func TestDefaultValueFloat(t *testing.T) {
    def := DefaultValue(3.14)
    loc := location.New(-112, -241)

    want := 3.14

    val, err := def(loc)
    if err != nil || val != want {
        t.Fatalf("def(%v) = %v, %v, want %v, %v", loc, val, err, want, nil)
    }
}


func TestWithDefault(t *testing.T) {
    want := 3

    g := WithDefault(want)
    loc := location.New(1,2)

    val, err := g.Get(loc)
    if val != want || err != nil {
        t.Fatalf("WithDefault(%v).Get(%v) = %v, %v, want %v, %v", want, loc, val, err, want, nil)
    }

    g.Set(loc, 13)

    val, err = g.Get(loc)
    if val != 13 || err != nil {
        t.Fatalf("WithDefault(%v).Get(%v) = %v, %v, want %v, %v (after Set)", want, loc, val, err, 13, nil)
    }
}

type WDF struct {
    Value string
    Error bool
}

func TestWithDefaultFunc(t *testing.T) {
    a := location.New(1, 2)
    b := location.New(2, 5)
    c := location.New(-1, 2)

    wants := map[location.Location]WDF{
        a:{"Marc", false},
        b:{"Mike", false},
        c: {"", true},
    }

    def := func(loc location.Location) (string, error) {
        if loc.X * loc.Y >= 0 {
            return "Mike", nil
        }
        return "", fmt.Errorf("NO NEGATIVE COORDINATES")
    }

    g := WithDefaultFunc(def)
    g.Set(a, "Marc")

    for loc, want := range wants {
        t.Run(fmt.Sprintf("%v:%v:%v", loc, want.Value, want.Error), func(t *testing.T) {
            val, err := g.Get(loc)

            if val != want.Value || (err != nil) != want.Error {
                t.Fatalf("WithDefaultFunc(...).Get(%v) = %q, %v, want %q, %v", loc, val, err != nil, want.Value, want.Error)
            }
        })
    }
}


func TestGetWithError(t *testing.T) {
    g := New[int]()
    loc := location.New(1,5)

    val, err := g.Get(loc)
    if val != 0 || err == nil {
        t.Fatalf("Get(%v) = %v, %v, want 0, %v", loc, val, err != nil, true)
    }
}

func TestGetKnown(t *testing.T) {
    g := New[string]()
    loc := location.New(2,-1)
    want := "hello, world!"

    g.Set(loc, want)

    val, err := g.Get(loc)
    if val != want || err != nil {
        t.Fatalf("Get(%v) = %q, %v, want %q, %v", loc, val, err != nil, want, false)
    }
}

func TestLen(t *testing.T) {
    g := New[string]()
    g.Set(location.New(1,2), "hello")

    length := g.Len()
    want := 1

    if length != want {
        t.Fatalf("g.Len() = %v, want %v", length, want)
    }

    g.Set(location.New(-1,34), "world")

    length = g.Len()
    want = 2

    if length != want {
        t.Fatalf("g.Len() = %v, want %v", length, want)
    }

    g.Remove(location.New(1,2))
    length = g.Len()
    want = 1

    if length != want {
        t.Fatalf("g.Len() = %v, want %v", length, want)
    }
}


func TestBounds(t *testing.T) {
    g := New[int]()

    tl := location.New(1,2)
    br := location.New(3,7)
    md := location.New(2,5)

    g.Set(tl, 3)
    g.Set(br, 4)
    g.Set(md, 8)

    want := Bounds{1,3,2,7}

    bounds, err := g.Bounds()

    if bounds != want || err != nil {
        t.Fatalf("g.Bounds() = %v, %v, want %v, %v", bounds, err != nil, want, false)
    }
}

func TestBoundsEmpty(t *testing.T) {
    g := New[int]()

    want := Bounds{0,0,0,0}

    bounds, err := g.Bounds()
    if bounds != want || err == nil {
        t.Fatalf("g.Bounds() = %v, %v, want %v, %v", bounds, err != nil, want, true)
    }
}

func contained(t *testing.T, b Bounds, loc location.Location, want bool) {
    if b.Contains(loc) != want {
        t.Fatalf("%v.Contains(%v) = %v, want %v", b, loc, !want, want)
    }
}

func TestBoundsContains(t *testing.T) {
    g := New[int]()

    topleft := location.New(-2,-2)
    bottomright := location.New(3,4)

    g.Set(topleft, 1)
    g.Set(bottomright, 1)

    bounds, _ := g.Bounds()

    outside := location.New(4,6)
    samex := location.New(1, 7)
    samey := location.New(4, 2)
    within := location.New(2,3)

    contained(t, bounds, within, true)
    contained(t, bounds, samex, false)
    contained(t, bounds, samey, false)
    contained(t, bounds, outside, false)

}
