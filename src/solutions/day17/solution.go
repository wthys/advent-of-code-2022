package day17

import (
    "fmt"
    "regexp"

    "github.com/sbwhitecap/tqdm"
    . "github.com/sbwhitecap/tqdm/iterators"

    "github.com/wthys/advent-of-code-2022/solver"
    "github.com/wthys/advent-of-code-2022/grid"
    "github.com/wthys/advent-of-code-2022/location"
    pf "github.com/wthys/advent-of-code-2022/pathfinding"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "17"
}

func (s solution) Part1(input []string) (string, error) {
    jetgen, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    g := prepareGrid()
    rockgen := NewRockGen()

    down := location.New(0,-1)

    maxHeight := 0

    limit := 2022
    cooldown := 100

    for i := 0; i < limit; i++ {
        rock := rockgen.Generate(location.New(2, maxHeight + 4))
        //fmt.Println()
        //printGridWithRock(g, rock)

        for {

            jet := jetgen.NextJet()

            blownRock := rock.Add(jet)
            if rockHitsWallOrDebris(g, blownRock) {
                blownRock = rock
            }

            fallenRock := blownRock.Add(down)
            if rockComesToRest(g, fallenRock) {
                landRock(g, blownRock)
                rock = blownRock
                break
            }

            rock = fallenRock
        }

        if rock.TopPiece().Y > maxHeight {
            //fmt.Printf("!! new max height -> %v (was %v) !!\n", rock.TopPiece().Y, maxHeight)
            maxHeight = rock.TopPiece().Y
        }

        cooldown -= 1
        if g.Len() > 1_000 && cooldown <= 0 {
            pruneGrid(g, maxHeight)
            cooldown = 100
        }

    }

    //fmt.Println("final state")
    //printGridWithRock(g, Rock{[]location.Location{{3, maxHeight+3}}})


    return solver.Solved(maxHeight)
}

func (s solution) Part2(input []string) (string, error) {
    jetgen, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    g := prepareGrid()
    rockgen := NewRockGen()

    down := location.New(0,-1)

    maxHeight := 0

    limit := 1_000_000_000_000
    cooldown := 1000

    tqdm.With(Interval(0, limit), "part 2 progress", func( _ interface{}) (brk bool) {
        rock := rockgen.Generate(location.New(2, maxHeight + 4))
        //fmt.Println()
        //printGridWithRock(g, rock)

        for {
            jet := jetgen.NextJet()

            blownRock := rock.Add(jet)
            if rockHitsWallOrDebris(g, blownRock) {
                blownRock = rock
            }

            fallenRock := blownRock.Add(down)
            if rockComesToRest(g, fallenRock) {
                landRock(g, blownRock)
                rock = blownRock
                break
            }

            rock = fallenRock
        }

        if rock.TopPiece().Y > maxHeight {
            maxHeight = rock.TopPiece().Y
        }

        if g.Len() > 1_000 && cooldown <= 0 {
            pruneGrid(g, maxHeight)
            cooldown = 100
        }
        cooldown -= 1

        return
    })

    //fmt.Println("final state")
    //printGridWithRock(g, Rock{[]location.Location{{3, maxHeight+3}}})


    return solver.Solved(maxHeight)
}

func landRock(g *grid.Grid[string], rock Rock) {
    rock.Do(func(loc location.Location) {
        g.Set(loc, "#")
    })
}

func pruneGrid(g *grid.Grid[string], maxHeight int) error {
    topleft := location.New(0,maxHeight)
    val := "."
    down := location.New(0,-1)

    val, _ = g.Get(topleft)
    for val != "#" {
        topleft = topleft.Add(down)
        val, _ = g.Get(topleft)
    }

    topright := location.New(6, maxHeight)
    val, _ = g.Get(topright)
    for val != "#" {
        topright = topright.Add(down)
        val, _ = g.Get(topright)
    }


    //fmt.Printf("searching route from %v -> %v    \r", topleft, topright)
    // find a route from left to right
    allLocs := []location.Location{}
    g.Apply(func(loc location.Location, value string) {
        if value == "#" {
            allLocs = append(allLocs, loc)
        }
    })

    neejbers := func(loc location.Location) []location.Location {
        nbh := []location.Location{}
        for x := -1; x <= 1; x++ {
            for y := -1; y <= 1; y++ {
                if x == 0 && y == 0 {
                    continue
                }

                pos := location.New(x, y).Add(loc)

                val, err := g.Get(pos)
                if err == nil && val == "#" {
                    nbh = append(nbh, pos)
                }
            }
        }
        return nbh
    }

    path, err := pf.ShortestPath(allLocs, topleft, topright, neejbers)
    if err != nil {
        return err
    }

    lowestPoint := maxHeight
    for _, loc := range path {
        if loc.Y < lowestPoint {
            lowestPoint = loc.Y
        }
    }

    //fmt.Printf("pruning all location below y=%v      \r", lowestPoint)

    n := 0
    g.Apply(func(loc location.Location, _ string) {
        if loc.Y < lowestPoint {
            g.Remove(loc)
            n += 1
        }
    })

    //fmt.Printf("pruned %v locations below y=%v      \n", n, lowestPoint)

    return nil
}

func rockComesToRest(g *grid.Grid[string], rock Rock) bool {
    hit := false
    rock.Do(func(loc location.Location) {
        if hit {
            return
        }

        if loc.Y <= 0 {
            hit = true
            return
        }

        val, err := g.Get(loc)
        if val == "#" || err != nil {
            hit = true
        }
    })

    return hit
}

func rockHitsWallOrDebris(g *grid.Grid[string], rock Rock) bool {
    hit := false
    rock.Do(func(loc location.Location) {
        if hit {
            return
        }

        val, err := g.Get(loc)
        if val == "#" || err != nil {
            hit = true
        }
    })
    return hit
}

func printGridWithRock(g *grid.Grid[string], rock Rock) {
    bounds := grid.Bounds{-1, 7, 1, rock.TopPiece().Y}


    for y := bounds.Ymax; y >= bounds.Ymin; y-- {
        for x := bounds.Xmin; x <= bounds.Xmax; x++ {
            pos := location.New(x, y)
            rockPrinted := false
            rock.Do(func(loc location.Location) {
                if loc == pos {
                    fmt.Print("@")
                    rockPrinted = true
                }
            })
            if rockPrinted {
                continue
            }

            val, err := g.Get(pos)
            if err != nil {
                val = "|"
            }
            fmt.Print(val)
        }
        fmt.Println()
    }
    fmt.Println("+-------+")

}

func prepareGrid() *grid.Grid[string] {
    maxWidthDefault := func(loc location.Location) (string, error) {
        if loc.X < 0 || loc.X >= 7 {
            return "", fmt.Errorf("hit the wall")
        }
        return ".", nil
    }

    return grid.WithDefaultFunc(maxWidthDefault)
}

func parseInput(input []string) (*JetGen, error) {
    reInput := regexp.MustCompile("^[<>]+$")

    if !reInput.MatchString(input[0]) {
        return nil, fmt.Errorf("input contains unexpected characters")
    }

    return NewJetGen(input[0]), nil
}

type (
    Rock struct {
        pieces []location.Location
    }

    RockGen struct {
        nextRock int
        rocks []Rock
    }

    JetGen struct {
        nextJet int
        jets string
    }
)

func NewJetGen(jets string) *JetGen {
    return &JetGen{0, jets}
}

func (g *JetGen) NextJet() location.Location {
    dir := g.jets[g.nextJet]
    g.nextJet = (g.nextJet + 1) % len(g.jets)
    switch dir {
    case '>':
        return location.New(1,0)
    default:
        return location.New(-1,0)
    }
}

func NewRockGen() *RockGen {
    rock1 := Rock{[]location.Location{
        {0,0}, {1,0}, {2,0}, {3,0},
    }}

    rock2 := Rock{[]location.Location{
               {1,2},
        {0,1}, {1,1}, {2,1},
               {1,0},
    }}

    rock3 := Rock{[]location.Location{
                      {2,2},
                      {2,1},
        {0,0}, {1,0}, {2,0},
    }}

    rock4 := Rock{[]location.Location{
        {0,3},
        {0,2},
        {0,1},
        {0,0},
    }}

    rock5 := Rock{[]location.Location{
        {0,1}, {1,1},
        {0,0}, {1,0},
    }}

    return &RockGen{0, []Rock{rock1, rock2, rock3, rock4, rock5}}
}

func (g *RockGen) Generate(loc location.Location) Rock {
    rock := g.rocks[g.nextRock].Add(loc)
    g.nextRock = (g.nextRock + 1) % len(g.rocks)
    return rock
}

func (r Rock) Add(loc location.Location) Rock {
    pieces := []location.Location{}
    for _, piece := range r.pieces {
        pieces = append(pieces, piece.Add(loc))
    }
    return Rock{pieces}
}

func (r Rock) Do(doer func(loc location.Location)) {
    for _, piece := range r.pieces {
        doer(piece)
    }
}

func (r Rock) TopPiece() location.Location {
    return r.pieces[0]
}
