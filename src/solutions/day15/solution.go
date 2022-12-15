package day15

import (
    "fmt"
    "sync"
    "time"
    "regexp"
    "strconv"

    "github.com/golang-collections/collections/set"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/grid"
    "github.com/wthys/advent-of-code-2022/location"
    "github.com/wthys/advent-of-code-2022/util"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "15"
}

func (s solution) Part1(input []string) (string, error) {
    pairs, err := parseInput(input)
    if err != nil {
        return "", err
    }


    poi := set.New()
    leftmost := 0
    rightmost := 0
    for _, pair := range pairs {
        poi.Insert(pair.Beacon)
        poi.Insert(pair.Sensor)

        left  := pair.Sensor.Add(location.New(-pair.Distance(), 0))
        right := pair.Sensor.Add(location.New(pair.Distance(), 0))

        if left.X < leftmost {
            leftmost = left.X
        }
        if right.X > rightmost {
            rightmost = right.X
        }
    }

    row := 2_000_000


    relevant_pairs := []Pair{}
    for _, pair := range pairs {
        loc := location.New(pair.Sensor.X, row)
        if pair.Contains(loc) {
            relevant_pairs = append(relevant_pairs, pair)
        }
    }

    fmt.Printf("%v pairs, %v are relevant y=%v\n", len(pairs), len(relevant_pairs), row)

    total := 0
    for x := leftmost; x <= rightmost; x++ {
        loc := location.New(x, row)

        // is it a beacon or sensor?
        if poi.Has(loc) {
            continue
        }

        // check if we're within range of a pair
        for _, pair := range relevant_pairs {
            if pair.Contains(loc) {
                total += 1
                // no need to check the rest
                break
            }
        }
    }


    return strconv.Itoa(total), nil
}

func (s solution) Part2(input []string) (string, error) {
    pairs, err := parseInput(input)
    if err != nil {
        return "", err
    }


    limit := 4_000_000
    //limit := 20


    center := location.New(0,0)
    for _, pair := range pairs {
        center = center.Add(pair.Sensor)
    }

    center.X = center.X / len(pairs)
    center.Y = center.Y / len(pairs)


    bounds := grid.Bounds{0,limit,0,limit}
    candidates := set.New()

    prog := Progress{}

    wg := sync.WaitGroup{}
    for offset := 1; offset <= 3; offset++{
        fmt.Printf("checking edge %v\n", offset)
        for _, pair := range pairs {
            pair.EdgeDo(offset, func(loc location.Location) {
                if true {
                    fmt.Printf("%v       \r",prog.Tick())
                }
                wg.Add(1)
                go func() {
                    if !bounds.Contains(loc) {
                        //fmt.Printf("%v out of bounds\n", loc)
                        return
                    }

                    for _, pair := range pairs {
                        if pair.Contains(loc) {
                            //fmt.Printf("%v claimed by %v\n", loc, pair.Sensor)
                            return
                        }
                    }

                    candidates.Insert(loc)
                    wg.Done()

                }()
            })
        }
    }

    wg.Wait()

    fmt.Printf("determining closest of %v candidates\n", candidates.Len())
    centermost := location.New(0,0)
    candidates.Do(func(el interface{}) {
        loc, _ := el.(location.Location)
        if ManhattanDistance(loc, center) < ManhattanDistance(centermost, center) {
            centermost = loc
        }
    })

    freq := TuningFrequency(centermost)

    fmt.Printf("%v => %v\n", centermost, freq)


    return strconv.Itoa(freq), nil
}

type (
    Pair struct {
        Sensor location.Location
        Beacon location.Location
    }

    Progress struct {
        idx int
        start *time.Time
    }
)

func (p *Progress) Tick() string {
    glyphs := "-/|\\"
    p.idx = (p.idx + 1) % len(glyphs)
    if p.start == nil {
        now := time.Now()
        p.start = &now
    }
    return fmt.Sprintf("  %v  %v", string(glyphs[p.idx]), time.Since(*(p.start)))
}

func prepareGrid(pairs []Pair) *grid.Grid[string] {
    g := grid.WithDefault(".")

    for _, pair := range pairs {
        g.Set(pair.Sensor, "S")
        g.Set(pair.Beacon, "B")
    }

    return g
}

func (p Pair) Distance() int {
    return ManhattanDistance(p.Sensor, p.Beacon)
}

func (p Pair) Contains(loc location.Location) bool {
    return ManhattanDistance(p.Sensor, loc) <= ManhattanDistance(p.Sensor, p.Beacon)
}

func (p Pair) EdgeDo(offset int, doer func(loc location.Location)) {
    dist := p.Distance() + offset
    for x := -dist - offset; x <= dist + offset; x++ {
        doer(location.New(x, x - dist).Add(p.Sensor))
        doer(location.New(x, dist - x).Add(p.Sensor))
    }
}

func (p Pair) Edge(offset int) *set.Set {
    edge := set.New()

    p.EdgeDo(offset, func(loc location.Location) {
        edge.Insert(loc)
    })

    return edge
}

func ManhattanDistance(a, b location.Location) int {
    diff := a.Subtract(b)
    return util.Abs(diff.X) + util.Abs(diff.Y)
}

func TuningFrequency(loc location.Location) int {
    return 4_000_000 * loc.X + loc.Y
}

func rotate90(loc location.Location) location.Location {
    return location.New(-loc.Y, loc.X)
}

func min(values ...int) int {
    least := values[0]
    for _, val := range values {
        if val < least {
            least = val
        }
    }
    return least
}

func SpiralDo(center location.Location, max int, doer func(loc location.Location) error) {
    loc := center
    motion := location.New(1, 0)
    extent := 0

    stop := make(chan int)

    wg := sync.WaitGroup{}

    count := max
    for count > 0 {
        if motion.Y == 0 {
            extent += 1
        }
        k := min(count, extent)
        for k > 0 {
            wg.Add(1)
            go func() {
                err := doer(loc)
                if err != nil {
                    stop <- 1
                }
                wg.Done()
            }()
            loc = loc.Add(motion)
            k -= 1
        }

        select {
        case <- stop:
            wg.Wait()
            return
        case <- time.After(time.Millisecond):
            wg.Wait()
        }
        motion = rotate90(motion)
    }
}

func parseInput(input []string) ([]Pair, error) {
    reSensor := regexp.MustCompile("x=(-?\\d+), y=(-?\\d+)")

    pairs := []Pair{}

    for n, line := range input {
        caps := reSensor.FindAllStringSubmatch(line, 2)
        if caps == nil {
            return nil, fmt.Errorf("could not make sense of line %v: %q", n, line)
        }

        xs, _ := strconv.Atoi(caps[0][1])
        ys, _ := strconv.Atoi(caps[0][2])
        xb, _ := strconv.Atoi(caps[1][1])
        yb, _ := strconv.Atoi(caps[1][2])

        pairs = append(pairs, Pair{location.New(xs, ys), location.New(xb, yb)})
    }

    return pairs, nil
}
