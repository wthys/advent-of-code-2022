package day6

import (
    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/collections/set"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "6"
}

const (
    PacketStart = 4
    MessageStart = 14
)

func isStartPacket(data string, index, length int) bool {
    check := set.New[rune]()
    for _, c := range data[index:index+length] {
        check.Add(c)
    }
    return check.Len() == length
}

func (s solution) Part1(input []string) (string, error) {

    data := input[0]

    for i := 0; i < len(data)-PacketStart; i++ {
        if isStartPacket(data, i, PacketStart) {
            return solver.Solved(i+PacketStart)
        }
    }

    return solver.Unsolved, nil
}

func (s solution) Part2(input []string) (string, error) {
    data := input[0]

    for i := 0; i < len(data)-MessageStart; i++ {
        if isStartPacket(data, i, MessageStart) {
            return solver.Solved(i+MessageStart)
        }
    }

    return solver.Unsolved, nil
}
