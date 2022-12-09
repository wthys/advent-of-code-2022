package util

import (
    "fmt"
    "testing"
)

type (
    testInt struct {
        Input int
        Want int
    }

    testFloat struct {
        Input float64
        Want float64
    }
)

func TestSignInt(t *testing.T) {
    cases := []testInt{
        {int(5), int(1)},
        {int(0), int(0)},
        {int(-245), int(-1)},
    }

    for _, cs := range cases {
        t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
            s := Sign(cs.Input)
            if s != cs.Want {
                t.Fatalf("Sign(%v) = %v, want %v", cs.Input, s, cs.Want)
            }
        })
    }
}

func TestSignFloat(t *testing.T) {
    cases := []testFloat{
        {float64(3.14), float64(1.0)},
        {float64(0.0), float64(0.0)},
        {float64(-5.256), float64(-1.0)},
    }

    for _, cs := range cases {
        t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
            s := Sign(cs.Input)
            if s != cs.Want {
                t.Fatalf("Sign(%v) = %v, want %v", cs.Input, s, cs.Want)
            }
        })
    }
}

func TestAbsInt(t *testing.T) {
    cases := []testInt{
        {int(5), int(5)},
        {int(0), int(0)},
        {int(-245), int(245)},
    }

    for _, cs := range cases {
        t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
            s := Abs(cs.Input)
            if s != cs.Want {
                t.Fatalf("Abs(%v) = %v, want %v", cs.Input, s, cs.Want)
            }
        })
    }
}

func TestAbsFloat(t *testing.T) {
    cases := []testFloat{
        {float64(3.14), float64(3.14)},
        {float64(0.0), float64(0.0)},
        {float64(-5.256), float64(5.256)},
    }

    for _, cs := range cases {
        t.Run(fmt.Sprintf("%v", cs.Input), func(t *testing.T) {
            s := Abs(cs.Input)
            if s != cs.Want {
                t.Fatalf("Abs(%v) = %v, want %v", cs.Input, s, cs.Want)
            }
        })
    }
}
