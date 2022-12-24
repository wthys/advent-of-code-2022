package day24

import (
    "fmt"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/location"
    "github.com/wthys/advent-of-code-2022/collections/set"
    pf "github.com/wthys/advent-of-code-2022/pathfinding"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "24"
}

func (s solution) Part1(input []string) (string, error) {
    valley, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    neejbers := func (state ValleyState) []ValleyState {
        return valley.LookIntoFuture(state)
    }

    //fmt.Println("   0: going the the exit!")
    hasExit := func (state ValleyState) bool {
        //dist := valley.Exit.Subtract(state.Expedition).Manhattan()
        //fmt.Printf("  time = %v, distance to target = %v     \r", state.Time, dist)
        return valley.ExpeditionAtExit(state)
    }
    d1 := pf.ControlledDijkstra(valley.State, neejbers, hasExit)
    exit := valley.State
    d1.DoNodes(func(state ValleyState) bool {
        if valley.ExpeditionAtExit(state) {
            exit = state
            return false
        }
        return true
    })
    //fmt.Println()

    return solver.Solved(exit.Time)
}

func (s solution) Part2(input []string) (string, error) {
    valley, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    neejbers := func (state ValleyState) []ValleyState {
        return valley.LookIntoFuture(state)
    }

    //fmt.Println("   0: going the the exit!")
    hasExit := func (state ValleyState) bool {
        //dist := valley.Exit.Subtract(state.Expedition).Manhattan()
        //fmt.Printf("  time = %v, distance to target = %v     \r", state.Time, dist)
        return valley.ExpeditionAtExit(state)
    }
    d1 := pf.ControlledDijkstra(valley.State, neejbers, hasExit)

    exit := valley.State
    d1.DoNodes(func(state ValleyState) bool {
        if valley.ExpeditionAtExit(state) {
            exit = state
            return false
        }
        return true
    })

    //fmt.Printf("%4v: seems we have to go back to the entrance... -_-\n", exit.Time)
    hasEntrance := func (state ValleyState) bool {
        //dist := valley.Entrance.Subtract(state.Expedition).Manhattan()
        //fmt.Printf("  time = %v, distance to target = %v     \r", state.Time, dist)
        return valley.ExpeditionAtEntrance(state)
    }
    d2 := pf.ControlledDijkstra(exit, neejbers, hasEntrance)

    entrance := exit
    d2.DoNodes(func(state ValleyState) bool {
        if valley.ExpeditionAtEntrance(state) {
            entrance = state
            return false
        }
        return true
    })

    //fmt.Printf("%4v: aaaaand back to the exit!\n", entrance.Time)
    d3 := pf.ControlledDijkstra(entrance, neejbers, hasExit)

    exit = entrance
    d3.DoNodes(func(state ValleyState) bool {
        if valley.ExpeditionAtExit(state) {
            exit = state
            return false
        }
        return true
    })
    //fmt.Println()

    return solver.Solved(exit.Time)
}

const (
    UP = "^"
    DOWN = "v"
    LEFT = "<"
    RIGHT = ">"
)

var (
    DIRECTIONS = map[string]location.Location{
        RIGHT: {1,0},
        LEFT: {-1,0},
        UP: {0,-1},
        DOWN: {0,1},
    }
)

type (
    Blizzard struct {
        Dir string
        Pos location.Location
    }

    BlizzardStateCache struct {
        times map[int][]Blizzard
    }

    ValleyState struct {
        Time int
        Expedition location.Location
    }

    Valley struct {
        Walls map[string]int
        Entrance location.Location
        Exit location.Location
        State ValleyState
        BSC BlizzardStateCache
    }
)

func (b *BlizzardStateCache) Store(time int, blizzards []Blizzard) {
    (*b).times[time] = blizzards
}

func (b BlizzardStateCache) Retrieve(time int) []Blizzard {
    blizzards, ok := b.times[time]
    if !ok {
        return nil
    }
    return blizzards
}

func (v Valley) Clamp(loc location.Location) location.Location {
    if loc == v.Entrance || loc == v.Exit {
        return loc
    }

    if loc.X < v.Walls[LEFT] {
        return location.New(v.Walls[RIGHT], loc.Y)
    }

    if loc.X > v.Walls[RIGHT] {
        return location.New(v.Walls[LEFT], loc.Y)
    }

    if loc.Y < v.Walls[UP] {
        return location.New(loc.X, v.Walls[DOWN])
    }

    if loc.Y > v.Walls[DOWN] {
        return location.New(loc.X, v.Walls[UP])
    }

    return loc
}

func (v Valley) ExpeditionAtExit(state ValleyState) bool {
    return state.Expedition == v.Exit
}

func (v Valley) ExpeditionAtEntrance(state ValleyState) bool {
    return state.Expedition == v.Entrance
}

func (v Valley) LookIntoFuture(state ValleyState) []ValleyState {
    candidates := set.New(state.Expedition)
    for _, dir := range DIRECTIONS {
        pos := state.Expedition.Add(dir)
        if pos == v.Exit || pos == v.Entrance || pos == v.Clamp(pos) {
            candidates.Add(pos)
        }
    }

    blizzards := v.BSC.Retrieve(state.Time + 1)
    if blizzards == nil {
        blizzards = []Blizzard{}
        for _, blizz := range v.BSC.Retrieve(state.Time) {
            pos := v.Clamp(blizz.Pos.Add(DIRECTIONS[blizz.Dir]))
            blizzards = append(blizzards, Blizzard{blizz.Dir, pos})
        }
        v.BSC.Store(state.Time + 1, blizzards)
    }
    for _, blizz := range blizzards {
        candidates.Remove(blizz.Pos)
    }

    states := []ValleyState{}
    candidates.Do(func(cand location.Location) bool {
        states = append(states, ValleyState{state.Time + 1, cand})
        return true
    })

    return states
}

func (v Valley) WithState(state ValleyState) Valley {
    return Valley{v.Walls, v.Entrance, v.Exit, state, v.BSC}
}

func (v Valley) String() string {
    str := strings.Builder{}

    poi := map[location.Location]string{}
    for _, blizz := range v.BSC.Retrieve(v.State.Time) {
        poi[blizz.Pos] += blizz.Dir
    }

    for y := v.Walls[UP]-1; y <= v.Walls[DOWN]+1; y++ {
        for x := v.Walls[LEFT]-1; x <= v.Walls[RIGHT]+1; x++ {
            pos := location.New(x,y)
            if pos == v.State.Expedition {
                fmt.Fprint(&str, "E")
                continue
            }
            if pos == v.Entrance || pos == v.Exit {
                fmt.Fprint(&str, ".")
                continue
            }
            if pos != v.Clamp(pos) {
                fmt.Fprint(&str, "#")
                continue
            }
            blizz, ok := poi[pos]
            if !ok {
                fmt.Fprint(&str, ".")
                continue
            }
            switch len(blizz) {
            case 0:
                fmt.Fprint(&str, ".")
            case 1:
                fmt.Fprint(&str, blizz)
            default:
                fmt.Fprintf(&str, "%v", len(blizz))
            }
        }
        fmt.Fprint(&str, "\n")
    }

    return str.String()
}

func parseInput(input []string) (Valley, error) {
    blizzards := []Blizzard{}
    entrance := location.New(0,0)
    exit := location.New(0,0)

    for y, line := range input {
        for x, ch := range line {
            if y == 0 {
                if string(ch) == "." {
                    entrance = location.New(x, y)
                }
                continue
            }

            if y == len(input)-1 {
                if string(ch) == "." {
                    exit = location.New(x, y)
                }
                continue
            }

            if string(ch) == "#"  {
                continue
            }

            if string(ch) == "." {
                continue
            }

            blizzards = append(blizzards, Blizzard{string(ch), location.New(x, y)})
        }
    }

    walls := map[string]int {
        UP: 1,
        DOWN: len(input)-2,
        LEFT: 1,
        RIGHT: len(input[0])-2,
    }

    cache := BlizzardStateCache{map[int][]Blizzard{}}
    cache.Store(0, blizzards)

    valley := Valley{walls, entrance, exit, ValleyState{0, entrance}, cache}
    return valley, nil
}
