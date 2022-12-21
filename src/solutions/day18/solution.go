package day18

import (
    "fmt"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/location"
    "github.com/wthys/advent-of-code-2022/collections/set"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "18"
}

func (s solution) Part1(input []string) (string, error) {
    points, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    grid := Grid3{}
    for _, point := range points {
        grid.Add(point)
    }

    total := 0
    grid.Do(func(_ location.Location3, value int) error {
        total += value
        return nil
    })

    return solver.Solved(total)
}

func (s solution) Part2(input []string) (string, error) {
    points, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    grid := Grid3{}
    xlo, xhi, ylo, yhi, zlo, zhi := 1000, 0, 1000, 0, 1000, 0
    for _, pt := range points {
        grid.Add(pt)

        if pt.X < xlo { xlo = pt.X }
        if pt.X > xhi { xhi = pt.X }

        if pt.Y < ylo { ylo = pt.Y }
        if pt.Y > yhi { yhi = pt.Y }

        if pt.Z < zlo { zlo = pt.Z }
        if pt.Z > zhi { zhi = pt.Z }
    }

    voids := Grid3{}
    for x := xlo-1; x <= xhi+1; x++ {
        for y := ylo-1; y <= yhi+1; y++ {
            for z := zlo-1; z <= zhi+1; z++ {
                pos := location.New3(x, y, z)
                _, ok := grid[pos]
                if ok {
                    continue
                }

                voids.Add(pos)
            }
        }
    }

    volumes := contiguousVolumes(&voids)
    for _, vol := range volumes {
        //fmt.Printf("volume #%v, size=%v: ", i, vol.Len())
        //printSet(vol, 6)

        outside := false
        vol.Do(func(loc location.Location3) bool {
            xborder := loc.X == xlo-1 || loc.X == xhi+1
            yborder := loc.Y == ylo-1 || loc.Y == yhi+1
            zborder := loc.Z == zlo-1 || loc.Z == zhi+1

            if xborder || yborder || zborder {
                outside = true
                return false
            }
            return true
        })

        if outside {
            //fmt.Println("  -> outside of the drop!")
            continue
        }

        vol.Do(func(loc location.Location3) bool {
            grid.Add(loc)
            return true
        })
    }

    total := 0
    grid.Do(func(_ location.Location3, sides int) error {
        total += sides
        return nil
    })

    return solver.Solved(total)
}

type (
    Grid3 map[location.Location3]int
)

func printSet[T comparable](s *set.Set[T], limit int) {
    n := 0
    more := false
    fmt.Print("{")
    s.Do(func(value T) bool {
        n += 1
        if n > limit && limit > 0 {
            more = true
            return true
        }
        fmt.Printf(" %v", value)
        return true
    })

    fmt.Print(" ")
    if more {
        fmt.Print("...")
    }
    fmt.Println("}")
}

func contiguousVolumes(grid *Grid3) []*set.Set[location.Location3] {
    volumes := []*set.Set[location.Location3]{}

    grid.Do(func(loc location.Location3, _ int) error {
        bordering := set.New[int]()
        for _, neejber := range neejbers3(loc) {
            for i, vol := range volumes {
                if vol.Has(neejber) {
                    bordering.Add(i)
                }
            }
        }

        switch bordering.Len() {
        case 0:
            volumes = append(volumes, set.New(loc))
        case 1:
            idx := bordering.Values()[0]
            volumes[idx].Add(loc)
        default:
            mergedVolume := set.New(loc)
            bordering.Do(func(idx int) bool {
                mergedVolume = mergedVolume.Union(volumes[idx])
                return true
            })

            newVolumes := []*set.Set[location.Location3]{}
            replaced := false
            for i, vol := range volumes {
                if bordering.Has(i) {
                    if replaced {
                        continue
                    }
                    newVolumes = append(newVolumes, mergedVolume)
                    replaced = true
                    continue
                }

                newVolumes = append(newVolumes, vol)
            }
            volumes = newVolumes
        }
        return nil
    })

    return volumes
}

func neejbers3(loc location.Location3) []location.Location3 {
    nbh := []location.Location3{
        {-1,0,0}, {1,0,0},
        {0,-1,0}, {0,1,0},
        {0,0,-1}, {0,0,1},
    }

    for i, offset := range nbh {
        nbh[i] = loc.Add(offset)
    }

    return nbh
}

func (g *Grid3) Add(loc location.Location3) {
    if g.Has(loc) {
        return
    }

    neejbers := 6

    for _, neejber := range neejbers3(loc) {
        val, ok := (*g)[neejber]
        if ok {
            (*g)[neejber] = val - 1
            neejbers -= 1
        }
    }

    (*g)[loc] = neejbers
}

func (g *Grid3) Has(loc location.Location3) bool {
    _, ok := (*g)[loc]
    return ok
}

func (g *Grid3) Do(doer func(loc location.Location3, value int) error) {
    for loc, value := range *g {
        err := doer(loc, value)
        if err != nil {
            return
        }
    }
}

func parseInput(input []string) ([]location.Location3, error) {
    locs := []location.Location3{}

    for _, line := range input {
        loc, err := location.FromString3(fmt.Sprintf("(%s)", line))
        if err != nil {
            return nil, err
        }
        locs = append(locs, loc)
    }

    return locs, nil
}
