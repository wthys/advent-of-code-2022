package day20

import (
    "testing"
)

type (
    Case struct {
        Orig []Item
        Idx int
        Want []Item
    }
)

func EqualItemList(actual []Item, want []Item, t *testing.T) {
    if len(actual) != len(want) {
        t.Fatalf("list has wrong length %v, want %v", len(actual), len(want))
    }

    for i := range want {
        if actual[i] != want[i] {
            t.Fatalf("got %v, want %v", actual, want)
        }
    }
}

func TestMix(t *testing.T) {
    cases := []Case{
        {
            []Item{ {0,0}, {1,1} },
            1,
            []Item{ {1,1}, {0,0} },
        },
        {
            []Item{ {1,1}, {2,2}, {3,0} },
            1,
            []Item{ {2,2}, {1,1}, {3,0} },
        },
        {
            []Item{ {3,0}, {1,1}, {2,2} },
            0,
            []Item{ {3,0}, {2,2}, {1,1} },
        },
        {
            []Item{ {6,0}, {1,1}, {2,2} },
            0,
            []Item{ {6,0}, {1,1}, {2,2} },
        },
        {
            []Item{ {4,0}, {1,1}, {2,2}, {3,3} },
            0,
            []Item{ {4,0}, {2,2}, {3,3}, {1,1} },
        },
        {
            []Item{ {8,0}, {1,1}, {2,2}, {3,3} },
            0,
            []Item{ {8,0}, {3,3}, {1,1}, {2,2} },
        },
        {
            []Item{ {144,0}, {1,1}, {2,2}, {3,3} },
            0,
            []Item{ {144,0}, {1,1}, {2,2}, {3,3} },
        },
        {
            []Item{ {-4,0}, {1,1}, {2,2}, {3,3} },
            0,
            []Item{ {-4,0}, {3,3}, {1,1}, {2,2} },
        },
        {
            []Item{ {-8,0}, {1,1}, {2,2}, {3,3} },
            0,
            []Item{ {-8,0}, {2,2}, {3,3}, {1,1} },
        },
        {
            []Item{ {-144,0}, {1,1}, {2,2}, {3,3} },
            0,
            []Item{ {-144,0}, {1,1}, {2,2}, {3,3} },
        },
    }

    for _, cs := range cases {
        EqualItemList(Mix(cs.Orig, cs.Idx), cs.Want, t)
    }
}
