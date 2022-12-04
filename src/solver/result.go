package solver

import (
    "fmt"
)

type Result struct{
    Name string
    Part1 string
    Part2 string
}

func (r Result) String() string {
    if r.Part1 == "" {
        r.Part1 = Unsolved
    }

    if r.Part2 == "" {
        r.Part2 = Unsolved
    }

    if r.Name == "" {
        r.Name = Unknown
    }

    return fmt.Sprintf("%v\t%v\t%v", r.Name, r.Part1, r.Part2)
}
