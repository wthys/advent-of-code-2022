package day22

import (
    "fmt"
    "strings"
    "strconv"
    "regexp"
    "math"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/grid"
    "github.com/wthys/advent-of-code-2022/location"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "22"
}

func (s solution) Part1(input []string) (string, error) {
    mover, instructions, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    for _, instr := range instructions {
        instr.Apply(mover)
    }

    scoring := map[location.Location]int{
        {1,0}: 0,
        {0,-1}: 3,
        {-1,0}: 2,
        {0,1}: 1,
    }

    //printGridMover(mover.(*GridMover))

    pos := mover.Position()
    dir := mover.Direction()
    score := 1000 * pos.Y + 4 * pos.X + scoring[dir]
    //fmt.Printf("scoring -> 1000 * %v + 4 * %v + %v%v = %v\n", pos.Y, pos.X, scoring[dir], dir, score)

    return solver.Solved(score)
}

func (s solution) Part2(input []string) (string, error) {
    mvr, instructions, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    gm := mvr.(*GridMover)

    mover := FromGridMover(gm)

    mover.Turn(0)

    //fmt.Println("== initial state ==")
    //printCubeMover(mover)

    for _, instr := range instructions {
        //fmt.Printf("== %v ==\n", instr)
        instr.Apply(mover)
        //printCubeMover(mover)
    }

    //printCubeMover(mover)

    scoring := map[location.Location]int{
        {1,0}: 0,
        {0,-1}: 3,
        {-1,0}: 2,
        {0,1}: 1,
    }

    pos := mover.Position()
    dir := mover.Direction()
    score := 1000 * pos.Y + 4 * pos.X + scoring[dir]
    //fmt.Printf("scoring -> 1000 * %v + 4 * %v + %v%v = %v\n", pos.Y, pos.X, scoring[dir], dir, score)

    return solver.Solved(score)
}

type (
    WrappingGrid struct{
        g *grid.Grid[string]
    }

    CubeWrappingGrid struct {
        g *grid.Grid[string]
    }

    Mover interface{
        Forward(dist int)
        Turn(dir int)
        Position() location.Location
        Direction() location.Location
    }

    GridMover struct {
        Pos location.Location
        Dir location.Location
        Grid *WrappingGrid
    }

    CubeMover struct {
        Pos location.Location
        Dir location.Location
        Grid *CubeWrappingGrid
    }

    Instr interface {
        Apply(Mover)
    }

    ForwardInstr struct {
        dist int
    }

    TurnInstr struct {
        dir int
    }
)

func FromGridMover(m *GridMover) *CubeMover {
    cg := &CubeWrappingGrid{(*m).Grid.g}
    return &CubeMover{m.Position(), m.Direction(), cg}
}

func (instr ForwardInstr) Apply(mover Mover) {
    mover.Forward(instr.dist)
}

func (instr ForwardInstr) String() string {
    return fmt.Sprintf("forward %v", instr.dist)
}

func (instr TurnInstr) Apply(mover Mover) {
    mover.Turn(instr.dir)
}

func (instr TurnInstr) String() string {
    mapp := map[int]string{1:"right",-1:"left"}
    return fmt.Sprintf("turn %v", mapp[instr.dir])
}

var (
    DIRS = map[location.Location]string {
        {1,0}: ">",
        {-1,0}: "<",
        {0,1}: "v",
        {0,-1}: "^",
    }
)

const (
    TOP = iota
    RIGHT
    BOTTOM
    LEFT
    FRONT
    BACK
)

func (m GridMover) Position() location.Location {
    return m.Pos
}

func (m GridMover) Direction() location.Location {
    return m.Dir
}

func (m *GridMover) Forward(dist int) {
    pos := (*m).Pos
    dir := (*m).Dir

    for i := 0; i < dist; i++ {
        newPos, val := (*m).Grid.WrappedGet(pos, dir)
        if val == "#" {
            break
        }
        (*m).Grid.g.Set(newPos, DIRS[dir])

        pos = newPos
    }
    (*m).Pos = pos
}

func (m *GridMover) Turn(dir int) {
    for dir < 0 { dir += 4 }
    for dir > 3 { dir -= 4 }

    mdir := (*m).Dir

    switch dir {
    case 1:
        (*m).Dir = location.New(-mdir.Y, mdir.X)
    case 2:
        (*m).Dir = location.New(-mdir.X, -mdir.Y)
    case 3:
        (*m).Dir = location.New(mdir.Y, -mdir.X)
    default:
        // nothing
    }

    (*m).Grid.g.Set(m.Position(), DIRS[m.Direction()])
}

func (m CubeMover) Position() location.Location {
    return m.Pos
}

func (m CubeMover) Direction() location.Location {
    return m.Dir
}

func (m *CubeMover) Forward(dist int) {
    pos := (*m).Pos
    dir := (*m).Dir

    for i := 0; i < dist; i++ {
        newPos, newDir, val := (*m).Grid.WrappedGet(pos, dir)
        if val == "#" {
            break
        }
        pos = newPos
        dir = newDir

        (*m).Grid.g.Set(pos, DIRS[dir])
    }
    (*m).Pos = pos
    (*m).Dir = dir
}

func (m *CubeMover) Turn(dir int) {
    for dir < 0 { dir += 4 }
    for dir > 3 { dir -= 4 }

    mdir := (*m).Dir

    switch dir {
    case 1:
        (*m).Dir = location.New(-mdir.Y, mdir.X)
    case 2:
        (*m).Dir = location.New(-mdir.X, -mdir.Y)
    case 3:
        (*m).Dir = location.New(mdir.Y, -mdir.X)
    default:
        // nothing
    }

    (*m).Grid.g.Set(m.Position(), DIRS[m.Direction()])
}

func printGridMover(m *GridMover) {
    b, _ := (*m).Grid.g.Bounds()
    for y := b.Ymin; y <= b.Ymax; y++ {
        for x := b.Xmin; x <= b.Xmax; x++ {
            pos := location.New(x, y)
            val, err := (*m).Grid.g.Get(pos)
            if err != nil {
                val = " "
            }
            fmt.Print(val)
        }
        fmt.Println()
    }
}

func printCubeMover(m *CubeMover) {
    b, _ := (*m).Grid.g.Bounds()
    for y := b.Ymin; y <= b.Ymax; y++ {
        for x := b.Xmin; x <= b.Xmax; x++ {
            pos := location.New(x, y)
            val, err := (*m).Grid.g.Get(pos)
            if err != nil {
                val = " "
            }
            fmt.Print(val)
        }
        fmt.Println()
    }
}

func (g WrappingGrid) WrappedGet(pos, dir location.Location) (location.Location, string) {
    newpos := pos.Add(dir)

    val, err := g.g.Get(newpos)
    if err == nil {
        return newpos, val
    }

    opdir := dir.Scale(-1)
    newpos = pos.Add(opdir)
    newval := ""
    for {
        srch := newpos.Add(opdir)
        val, err := g.g.Get(srch)
        if err != nil {
            break
        }
        newpos = srch
        newval = val
    }

    return newpos, newval
}

func (g CubeWrappingGrid) nextLocDir(start, dir location.Location) (location.Location, location.Location) {
    area := g.g.Len() / 6
    size := 1
    if area == 2500 {
        size = 50
    } else {
        size = 4
    }

    sx := int(math.Floor(float64(start.X - 1) / float64(size)))
    sy := int(math.Floor(float64(start.Y - 1) / float64(size)))
    rstart := location.New(sx, sy)

    end := start.Add(dir)
    ex := int(math.Floor(float64(end.X - 1) / float64(size)))
    ey := int(math.Floor(float64(end.Y - 1) / float64(size)))
    rend := location.New(ex, ey)

    mapping := map[int]map[location.Location]int {
        50: {
            {1,0}: TOP,
            {2,0}: RIGHT,
            {1,1}: FRONT,
            {0,2}: LEFT,
            {1,2}: BOTTOM,
            {0,3}: BACK,
        },
        4: {
            {2,0}: TOP,
            {0,1}: BACK,
            {1,1}: LEFT,
            {2,1}: FRONT,
            {2,2}: BOTTOM,
            {3,2}: RIGHT,
        },
    }

    sface, _ := mapping[size][rstart]
    _, ok := mapping[size][rend]

    //fmt.Printf("%v/%v -> %v/%v ? %v\n", start, rstart, end, rend, ok)

    if ok {
        return end, dir
    }

    DIRECTIONS := map[int]location.Location{
        TOP: location.New(0,-1),
        LEFT: location.New(-1,0),
        BOTTOM: location.New(0,1),
        RIGHT: location.New(1,0),
    }

    dirmap := map[int]map[int]func(loc location.Location) (newloc, newdir location.Location){
        50: {
            /*     0  1  2
             *       +--+--+
             *  0    |T |R |
             *       +--+--+
             *  1    |F |
             *    +--+--+
             *  2 |L |BO|
             *    +--+--+
             *  3 |BA|
             *    +--+
             *
             */
            TOP: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[TOP] {
                    return location.New(1, 100 + loc.X), DIRECTIONS[RIGHT]
                }

                return location.New(1, 151 - loc.Y), DIRECTIONS[RIGHT]
            },
            FRONT: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[LEFT] {
                    return location.New(loc.Y - 50, 101), DIRECTIONS[BOTTOM]
                }

                return location.New(loc.Y + 50, 50), DIRECTIONS[TOP]
            },
            RIGHT: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[TOP] {
                    return location.New(loc.X - 100, 200), DIRECTIONS[TOP]
                }
                if dir == DIRECTIONS[RIGHT] {
                    return location.New(100, 151 - loc.Y), DIRECTIONS[LEFT]
                }

                return location.New(100, loc.X - 50), DIRECTIONS[LEFT]
            },
            BOTTOM: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[RIGHT] {
                    return location.New(150, 151 - loc.Y), DIRECTIONS[LEFT]
                }

                return location.New(50, loc.X + 100), DIRECTIONS[LEFT]
            },
            LEFT: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[TOP] {
                    return location.New(51, 50 + loc.X), DIRECTIONS[RIGHT]
                }

                return location.New(51, 151 - loc.Y), DIRECTIONS[RIGHT]
            },
            BACK: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[LEFT] {
                    return location.New(loc.Y - 100, 1), DIRECTIONS[BOTTOM]
                }

                if dir == DIRECTIONS[BOTTOM] {
                    return location.New(loc.X + 100, 1), DIRECTIONS[BOTTOM]
                }

                return location.New(loc.Y - 100, 150), DIRECTIONS[TOP]
            },
        },
        4: {
            /*
             *          +--+
             *          |T |
             *    +--+--+--+
             *    |BA|L |F |
             *    +--+--+--+--+
             *          |BO|R |
             *          +--+--+
             *
             */
            TOP: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[TOP] {
                    diff := loc.X - 12
                    return location.New(1 - diff, 5), DIRECTIONS[BOTTOM]
                }

                if dir == DIRECTIONS[LEFT] {
                    diff := loc.Y - 1
                    return location.New(5 + diff, 5), DIRECTIONS[BOTTOM]
                }

                diff := loc.Y - 1
                return location.New(16, 11 - diff), DIRECTIONS[LEFT]
            },
            FRONT: func(loc location.Location) (newloc, newdir location.Location) {
                diff := loc.Y - 5
                return location.New(16 - diff, 9), DIRECTIONS[BOTTOM]
            },
            RIGHT: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[TOP] {
                    diff := loc.X - 16
                    return location.New(12, 5 - diff), DIRECTIONS[LEFT]
                }

                if dir == DIRECTIONS[RIGHT] {
                    diff := loc.Y - 9
                    return location.New(12, 4 - diff), DIRECTIONS[LEFT]
                }

                diff := loc.X - 13
                return location.New(1, 8 - diff), DIRECTIONS[RIGHT]
            },
            BOTTOM: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[BOTTOM] {
                    diff := loc.X - 9
                    return location.New(4 - diff, 8), DIRECTIONS[TOP]
                }

                diff := loc.Y - 9
                return location.New(8 - diff, 8), DIRECTIONS[TOP]
            },
            LEFT: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[TOP] {
                    diff := loc.X - 5
                    return location.New(9, 1 + diff), DIRECTIONS[RIGHT]
                }

                diff := loc.X - 5
                return location.New(9, 12 - diff), DIRECTIONS[RIGHT]
            },
            BACK: func(loc location.Location) (newloc, newdir location.Location) {
                if dir == DIRECTIONS[LEFT] {
                    diff := loc.Y - 8
                    return location.New(13 - diff, 12), DIRECTIONS[TOP]
                }

                if dir == DIRECTIONS[BOTTOM] {
                    diff := loc.X - 4
                    return location.New(9 - diff, 12), DIRECTIONS[TOP]
                }

                diff := loc.X - 4
                return location.New(9 - diff, 1), DIRECTIONS[BOTTOM]
            },
        },
    }

    pp, dd := dirmap[size][sface](end)

    //fmt.Printf("going from (%v,%v) %v/%v to %v/%v\n", sface, DIRS[dir], end, dir, pp, dd)

    return pp, dd

}

func (g CubeWrappingGrid) WrappedGet(pos, dir location.Location) (location.Location, location.Location, string) {
    newpos, newdir := g.nextLocDir(pos, dir)

    val, err := g.g.Get(newpos)
    if err == nil {
        return newpos, newdir, val
    }

    fmt.Printf("|->  %v/%v -> %v/%v\n", pos, dir, newpos, newdir)
    panic(err)

}

func parseInput(input []string) (Mover, []Instr, error) {
    g := WrappingGrid{ grid.New[string]() }

    reInstr := regexp.MustCompile("(\\d+|[LR])")
    instructions := []Instr{}

    gridDone := false
    for y, line := range input {
        if len(strings.TrimSpace(line)) == 0 {
            gridDone = true
            continue
        }

        if gridDone {
            caps := reInstr.FindAllString(line, -1)
            if caps == nil {
                return nil, nil, fmt.Errorf("invalid instructions")
            }

            for _, instr := range caps {
                fwd, err := strconv.Atoi(instr)
                if err != nil {
                    switch instr {
                    case "L":
                        instructions = append(instructions, TurnInstr{-1})
                    case "R":
                        instructions = append(instructions, TurnInstr{1})
                    default:
                        return nil, nil, fmt.Errorf("could not make sense of %q", instr)
                    }
                    continue
                }

                instructions = append(instructions, ForwardInstr{fwd})
            }
        } else {
            for x, val := range line {
                tile := string(val)
                if tile != "." && tile != "#" {
                    continue
                }
                g.g.Set(location.New(x + 1, y + 1), tile)
            }
        }
    }

    bounds, _ := g.g.Bounds()
    start := location.New(0,0)
    for x := bounds.Xmin; x <= bounds.Xmax; x++ {
        start = location.New(x, bounds.Ymin)
        _, err := g.g.Get(start)
        if err == nil {
            break
        }
    }

    mover := GridMover{start, location.New(1,0), &g}

    return &mover, instructions, nil

}
