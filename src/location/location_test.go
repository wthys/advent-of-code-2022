package location

import (
    "fmt"
    "testing"
)

func TestNew(t *testing.T) {
    for x := -100; x <= 100 ; x++ {
        for y := -100; y <= 100; y++ {
            loc := New(x, y)
            want := Location{x, y}
            if want != loc {
                t.Fatalf(`New(%v, %v) = %v, want %v`, x, y, loc, want)
            }
        }
    }
}

type testFS struct {
    Input string
    Want Location
    Error bool
}

func TestFromString(t *testing.T) {
    cases := []testFS{
        {"(1,2)", New(1,2), false},
        {"( 2,  4  )", New(2,4), false},
        {"   ( 3    ,    9    )    ", New(3,9), false},
        {"(  1035, -19  )", New(1035,-19), false},
        {" ( hello, world) ", Location{}, true},
        {" some outrageous input 12, 45", Location{}, true},
        {" another outrageous input (12, 45)", Location{}, true},
    }

    for _, cs := range cases {
        t.Run(cs.Input, func(t *testing.T) {
            loc, err := FromString(cs.Input)
            if (cs.Error && err == nil) || (!cs.Error && err != nil) || cs.Want != loc {
                e := err != nil
                t.Fatalf("FromString(%q) = (%v, %v), want (%v, %v)", cs.Input, loc, e, cs.Want, cs.Error)
            }
        })
    }
}

func TestString(t *testing.T) {
    loc := New(5, 7)
    want := "(5,7)"

    if want != loc.String() {
        t.Fatalf("Location{5,7} = %q, want %q", loc, want)
    }
}

func TestAdd(t *testing.T) {
    a := New(1,2)
    b := New(3,4)

    want := New(4,6)
    c := a.Add(b)

    if want != c {
        t.Fatalf("%v + %v = %v, want %v", a, b, c, want)
    }

    d := b.Add(a)
    if want != d {
        t.Fatalf("%v + %v = %v, want %v", b, a, d, want)
    }
}

func TestSubtract(t * testing.T) {
    a := New(1,2)
    b := New(3,4)

    want := New(-2,-2)

    c := a.Subtract(b)
    if want != c {
        t.Fatalf("%v - %v = %v, want %v", a, b, c, want)
    }

    d := b.Subtract(a)
    want = New(2,2)
    if want != d {
        t.Fatalf("%v - %v = %v, want %v", b, a, d, want)
    }
}

func TestScale(t *testing.T) {
    a := New(2,7)

    want := New(-2,-7)

    s := a.Scale(-1)
    if s != want {
        t.Fatalf("%v.Scale(-1) = %v, want %v", a, s, want)
    }

    want = New(6, 21)
    s = a.Scale(3)
    if s != want {
        t.Fatalf("%v.Scale(3) = %v, want %v", a, s, want)
    }
}

type testU struct {
    loc Location
    want Location
}

func TestUnit(t *testing.T) {
    cases := []testU{
        {New(3,5), New(1,1)},
        {New(0,1), New(0,1)},
        {New(-1,198), New(-1,1)},
        {New(-1243,0), New(-1,0)},
        {New(-12,-56), New(-1,-1)},
        {New(0,0), New(0,0)},
        {New(1,1), New(1,1)},
    }

    for _, cs := range cases {
        t.Run(fmt.Sprintf("Unit(%v)", cs.loc), func (t *testing.T) {
            unit := cs.loc.Unit()
            if unit != cs.want {
                t.Fatalf("%v.Unit() = %v, want %v", cs.loc, unit, cs.want)
            }
        })
    }
}
