package solver

import (
    "fmt"
    "time"
)

type Result struct{
    Name string
    Part1 string
    Part2 string
    Elapsed []time.Duration
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

    if r.Elapsed != nil && len(r.Elapsed) == 2 {
        return fmt.Sprintf("%v\t%v\t%v\t%v\t%v", r.Name, r.Part1, r.Part2, r.Elapsed[0], r.Elapsed[1])
    }

    return fmt.Sprintf("%v\t%v\t%v", r.Name, r.Part1, r.Part2)
}
