package day15

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/wthys/advent-of-code-2022/collections/set"
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


    poi := set.New[location.Location]()
    leftmost := 0
    rightmost := 0
    for _, pair := range pairs {
        poi.Add(pair.Beacon)
        poi.Add(pair.Sensor)

        left  := pair.Sensor.Add(location.New(-pair.Distance(), 0))
        right := pair.Sensor.Add(location.New(pair.Distance(), 0))

        if left.X < leftmost {
            leftmost = left.X
        }
        if right.X > rightmost {
            rightmost = right.X
        }
    }

    //row := 10
    row := 2_000_000

    relevant_pairs := []Pair{}
    for _, pair := range pairs {
        loc := location.New(pair.Sensor.X, row)
        if pair.Contains(loc) {
            relevant_pairs = append(relevant_pairs, pair)
        }
    }

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


    //limit := 20
    limit := 4_000_000

    bounds := grid.Bounds{0,limit,0,limit}

    found := false

    distressbeacon := location.New(0,0)

    for offset := 1; offset <= 3 && !found; offset++{
        //fmt.Printf("checking edge %v\n", offset)
        for _, pair := range pairs {
            pair.EdgeDo(offset, func(loc location.Location) {
                if found {
                    return
                }

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

                found = true
                distressbeacon = loc
            })
        }
    }

    freq := TuningFrequency(distressbeacon)
    return strconv.Itoa(freq), nil
}

type (
    Pair struct {
        Sensor location.Location
        Beacon location.Location
    }
)

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

func ManhattanDistance(a, b location.Location) int {
    diff := a.Subtract(b)
    return util.Abs(diff.X) + util.Abs(diff.Y)
}

func TuningFrequency(loc location.Location) int {
    return 4_000_000 * loc.X + loc.Y
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
