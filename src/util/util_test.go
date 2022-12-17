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

func hash(array []int) int {
    h := 0
    for _, value := range array {
        h = 13*h + value
    }
    return h
}

func TestPermutationDo(t *testing.T) {
    array := []int{1,2,3}

    check := map[int]bool{
        hash([]int{1,2,3}): false,
        hash([]int{1,3,2}): false,
        hash([]int{2,1,3}): false,
        hash([]int{2,3,1}): false,
        hash([]int{3,2,1}): false,
        hash([]int{3,1,2}): false,
    }

    PermutationDo(3, array, func(perm []int) {
        check[hash(perm)] = true
    })

    for value, seen := range check {
        if !seen {
            t.Fatalf("PermutationDo(3, %v, ...) should produce %v, but was not seen", array, value)
        }
    }

}
