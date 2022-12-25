package day19

import (
    "fmt"
    "sync"
    "regexp"
    "strconv"
    "strings"

    "github.com/wthys/advent-of-code-2022/solver"
    pf "github.com/wthys/advent-of-code-2022/pathfinding"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "19"
}

func (s solution) Part1(input []string) (string, error) {
    blueprints, err := parseInput(input)
    if err != nil {
        return solver.Error(err)
    }

    neejbers := func (tick BPTick) []BPTick {
        if tick.tick >= 24 {
            return []BPTick{}
        }
        neejbers := []BPTick{}
        next := tick.tick + 1

        orecost := tick.bp.OreBotCost
        claycost := tick.bp.ClayBotCost
        obscost := tick.bp.ObsidianBotCost
        geocost := tick.bp.GeodeBotCost

        nOre := tick.bp.nOreBots
        nClay := tick.bp.nClayBots
        nObs := tick.bp.nObsidianBots

        oreNeeded := max(orecost.Ore, claycost.Ore, obscost.Ore, geocost.Ore)
        clayNeeded := max(orecost.Clay, claycost.Clay, obscost.Clay, geocost.Clay)
        obsNeeded := max(orecost.Obsidian, claycost.Obsidian, obscost.Obsidian, geocost.Obsidian)

        inv := tick.bp.Inventory
        /*
        timeLeft := 24 - tick.tick - 1

        needOre := (timeLeft * nOre + inv.Ore) < (timeLeft * oreNeeded)
        needClay := (timeLeft * nClay + inv.Clay) < (timeLeft * clayNeeded)
        needObs := (timeLeft * nObs + inv.Obsidian) < (timeLeft * obsNeeded)

        //*/

        geoBP := tick.bp
        if geoBP.PayGeodeBot() {
            geoBP.Produce()
            geoBP.AddGeodeBot()
            neejbers = append(neejbers, BPTick{geoBP, next})
            return neejbers
        }

        if nObs < obsNeeded {
            obsBP := tick.bp
            if obsBP.PayObsidianBot() {
                obsBP.Produce()
                obsBP.AddObsidianBot()
                neejbers = append(neejbers, BPTick{obsBP, next})
            }
        }

        if nClay < clayNeeded {
            clayBP := tick.bp
            if clayBP.PayClayBot() {
                clayBP.Produce()
                clayBP.AddClayBot()
                neejbers = append(neejbers, BPTick{clayBP, next})
            }
        }

        if nOre < oreNeeded {
            oreBP := tick.bp
            if oreBP.PayOreBot() {
                oreBP.Produce()
                oreBP.AddOreBot()
                neejbers = append(neejbers, BPTick{oreBP, next})
            }
        }

        if inv.LessThanOrEqual(geocost) || inv.LessThanOrEqual(obscost) || inv.LessThanOrEqual(claycost) || inv.LessThanOrEqual(orecost) {
            prodBP := tick.bp
            prodBP.Produce()
            neejbers = append(neejbers, BPTick{prodBP, next})
        }


        return neejbers
    }


    watcher := func(node BPTick) bool {
        inv := node.bp.Inventory
        fmt.Printf("BP %2v best @ %v -> %3v, %3v, %3v, %3v\n", node.bp.Id, node.tick, inv.Geode, inv.Obsidian, inv.Clay, inv.Ore)
        return false
    }



    totalQuality := 0
    totalMutex := sync.Mutex{}
    wg := sync.WaitGroup{}
    for _, bp := range blueprints {
        BP := *bp
        fmt.Printf("started processing bp %v\n", BP.Id)

        fmt.Printf("  bp %2v: ore cost = %v\n", BP.Id, BP.OreBotCost)
        fmt.Printf("  bp %2v: clay cost = %v\n", BP.Id, BP.ClayBotCost)
        fmt.Printf("  bp %2v: obsidian cost = %v\n", BP.Id, BP.ObsidianBotCost)
        fmt.Printf("  bp %2v: geode cost = %v\n", BP.Id, BP.GeodeBotCost)

        wg.Add(1)
        go func() {
            defer wg.Done()

            start := BPTick{BP, 0}
            d := pf.ControlledDijkstra(start, neejbers, watcher)

            bestBP := BP
            best := 0
            d.DoNodes(func(node BPTick) bool {
                if node.tick != 24 {
                    return true
                }

                if node.bp.Inventory.Geode > best {
                    bestBP = node.bp
                    best = node.bp.Inventory.Geode
                }
                return true
            })

            quality := bestBP.Id * best
            fmt.Printf("bp %v has %v geodes -> %v\n", bestBP.Id, best, quality)

            totalMutex.Lock()
            totalQuality += quality
            totalMutex.Unlock()
        }()
    }
    wg.Wait()


    return solver.Solved(totalQuality)
}

func (s solution) Part2(input []string) (string, error) {
    return solver.NotImplemented()
}

type (
    BotCost struct {
        Ore int
        Clay int
        Obsidian int
        Geode int
    }

    Blueprint struct {
        Id int
        OreBotCost BotCost
        ClayBotCost BotCost
        ObsidianBotCost BotCost
        GeodeBotCost BotCost
        Inventory BotCost
        nOreBots int
        nClayBots int
        nObsidianBots int
        nGeodeBots int
    }

    BPTick struct {
        bp Blueprint
        tick int
    }
)

func max(values ...int) int {
    best := values[0]

    for _, val := range values {
        if val > best {
            best = val
        }
    }
    return best
}

func min(values ...int) int {
    best := values[0]

    for _, val := range values {
        if val < best {
            best = val
        }
    }
    return best
}

func EnoughMaterial(bots int, inventory int, cost int, timeLeft int) int {
    // enough production to pay for all future bots
    if bots >= cost {
        return timeLeft
    }

    // amount of bots that can be built in the time left if production
    // stays unchanged 
    //
    // inv + (tl - 1) * bots >= tl * cost
    //   => (inv - bots) / (cost - bots) >= tl
    potentialBots := (inventory - bots) / (cost - bots)
    return max(0, potentialBots)
}

func EnoughForGeodeBots(bp Blueprint, timeLeft int) int {
    geobot := bp.GeodeBotCost
    enoughOre := EnoughMaterial(bp.nOreBots, bp.Inventory.Ore, geobot.Ore, timeLeft)
    enoughObs := EnoughMaterial(bp.nObsidianBots, bp.Inventory.Obsidian, geobot.Obsidian, timeLeft)
    return min(enoughOre, enoughObs, timeLeft)
}

func EnoughObsBots(bp Blueprint) bool {
    geobot := bp.GeodeBotCost

    needed := max(0, geobot.Obsidian - bp.nObsidianBots)

    return needed == 0
}

func EnoughClayBots(bp Blueprint) bool {
    obsbot := bp.ObsidianBotCost

    needed := max(0, obsbot.Clay - bp.nClayBots)

    return needed == 0
}

func EnoughOreBots(bp Blueprint) bool {
    geobot := bp.GeodeBotCost
    obsbot := bp.ObsidianBotCost
    claybot :=bp.ClayBotCost

    needed := max(0, obsbot.Ore - bp.nOreBots, geobot.Ore - bp.nOreBots, claybot.Ore - bp.nOreBots)

    return needed == 0
}


func MaximizeBlueprint(blueprint Blueprint, timeLeft int) Blueprint {
    if timeLeft <= 0 {
        return blueprint
    }

    //fmt.Printf("maximizing %v @ %v     \r", blueprint, timeLeft)

    bestBP := blueprint
    //bestMutex := sync.Mutex{}
    //wg := sync.WaitGroup{}

    geoBP := blueprint
    if geoBP.PayGeodeBot() {
        geoBP.Produce()
        geoBP.AddGeodeBot()

        production := EnoughForGeodeBots(geoBP, timeLeft - 1)
        for i := 0; i < production; i++ {
            geoBP.PayGeodeBot()
            geoBP.Produce()
            geoBP.AddGeodeBot()
        }

        //wg.Add(1)
        //go func() {
        //    defer wg.Done()
            geoBP = MaximizeBlueprint(geoBP, timeLeft - 1 - production)

        //    bestMutex.Lock()
            if bestBP.Inventory.LessThanOrEqual(geoBP.Inventory) {
                bestBP = geoBP
            }
        //    bestMutex.Unlock()
        //}()
    }

    if !EnoughObsBots(blueprint) {
        obsBP := blueprint
        if obsBP.PayObsidianBot() {
            obsBP.Produce()
            obsBP.AddObsidianBot()

            //wg.Add(1)
            //go func() {
            //    defer wg.Done()
                obsBP = MaximizeBlueprint(obsBP, timeLeft - 1)

            //    bestMutex.Lock()
                if bestBP.Inventory.LessThanOrEqual(obsBP.Inventory) {
                    bestBP = obsBP
                }
            //    bestMutex.Unlock()
            //}()
        }
    }

    if !EnoughClayBots(blueprint) {
        clayBP := blueprint
        if clayBP.PayClayBot() {
            clayBP.Produce()
            clayBP.AddClayBot()

            //wg.Add(1)
            //go func() {
            //    defer wg.Done()
                clayBP = MaximizeBlueprint(clayBP, timeLeft - 1)

            //    bestMutex.Lock()
                if bestBP.Inventory.LessThanOrEqual(clayBP.Inventory) {
                    bestBP = clayBP
                }
            //    bestMutex.Unlock()
            //}()
        }
    }

    if !EnoughOreBots(blueprint) {
        oreBP := blueprint
        if oreBP.PayOreBot() {
            oreBP.Produce()
            oreBP.AddOreBot()

            //wg.Add(1)
            //go func() {
            //    defer wg.Done()
                oreBP = MaximizeBlueprint(oreBP, timeLeft - 1)

            //    bestMutex.Lock()
                if bestBP.Inventory.LessThanOrEqual(oreBP.Inventory) {
                    bestBP = oreBP
                }
            //    bestMutex.Unlock()
            //}()
        }
    }

    prodBP := blueprint
    prodBP.Produce()
    //wg.Add(1)
    //go func() {
    //    defer wg.Done()
        prodBP = MaximizeBlueprint(prodBP, timeLeft - 1)

    //    bestMutex.Lock()
        if bestBP.Inventory.LessThanOrEqual(prodBP.Inventory) {
            bestBP = prodBP
        }
    //    bestMutex.Unlock()
    //}()

    //wg.Wait()


    fmt.Printf("current best blueprint %v -> %v quality      \r", bestBP, bestBP.Id * bestBP.Inventory.Geode)

    return bestBP
}

func (bc BotCost) Subtract(o BotCost) BotCost {
    return BotCost{
        Ore: bc.Ore - o.Ore,
        Clay: bc.Clay - o.Clay,
        Obsidian: bc.Obsidian - o.Obsidian,
        Geode: bc.Geode - o.Geode,
    }
}

func (bc BotCost) LessThanOrEqual(o BotCost) bool {
    if bc.Geode > o.Geode {
        return false
    }

    if bc.Obsidian > o.Obsidian {
        return false
    }

    if bc.Clay > o.Clay {
        return false
    }

    return bc.Ore <= o.Ore
}


func (bp *Blueprint) Produce() {
    (*bp).Inventory.Ore      += (*bp).nOreBots
    (*bp).Inventory.Clay     += (*bp).nClayBots
    (*bp).Inventory.Obsidian += (*bp).nObsidianBots
    (*bp).Inventory.Geode    += (*bp).nGeodeBots
}

func (bp *Blueprint) PayOreBot() bool {
    if !(*bp).OreBotCost.LessThanOrEqual((*bp).Inventory) {
        return false
    }

    (*bp).Inventory = (*bp).Inventory.Subtract((*bp).OreBotCost)
    return true
}

func (bp *Blueprint) AddOreBot() {
    (*bp).nOreBots += 1
}

func (bp *Blueprint) PayClayBot() bool {
    if !(*bp).ClayBotCost.LessThanOrEqual((*bp).Inventory) {
        return false
    }

    (*bp).Inventory = (*bp).Inventory.Subtract((*bp).ClayBotCost)
    return true
}

func (bp *Blueprint) AddClayBot() {
    (*bp).nClayBots += 1
}

func (bp *Blueprint) PayObsidianBot() bool {
    if !(*bp).ObsidianBotCost.LessThanOrEqual((*bp).Inventory) {
        return false
    }

    (*bp).Inventory = (*bp).Inventory.Subtract((*bp).ObsidianBotCost)
    return true
}

func (bp *Blueprint) AddObsidianBot() {
    (*bp).nObsidianBots += 1
}

func (bp *Blueprint) PayGeodeBot() bool {
    if !(*bp).GeodeBotCost.LessThanOrEqual((*bp).Inventory) {
        return false
    }

    (*bp).Inventory = (*bp).Inventory.Subtract((*bp).GeodeBotCost)
    return true
}

func (bp *Blueprint) AddGeodeBot() {
    (*bp).nGeodeBots += 1
}

func (bp Blueprint) String() string {
    str := strings.Builder{}
    fmt.Fprint(&str, "{")
    fmt.Fprintf(&str, " Id: %v", bp.Id)
    inv := bp.Inventory
    fmt.Fprintf(&str, ", Inv: { %v, %v, %v, %v }", inv.Ore, inv.Clay, inv.Obsidian, inv.Geode)
    fmt.Fprintf(&str, ", Bots: { %v, %v, %v, %v }", bp.nOreBots, bp.nClayBots, bp.nObsidianBots, bp.nGeodeBots)
    fmt.Fprint(&str, " }")
    return str.String()
}


func parseInput(input []string) ([]*Blueprint, error) {
    blueprints := []*Blueprint{}

    reNum := regexp.MustCompile("-?\\d+")

    for nr, line := range input {
        caps := reNum.FindAllString(line, -1)
        if caps == nil || len(caps) < 7 {
            return nil, fmt.Errorf("could not parse #%v: %v", nr, line)
        }

        id, _ := strconv.Atoi(caps[0])

        orebot := BotCost{}
        orebot.Ore, _ = strconv.Atoi(caps[1])

        claybot := BotCost{}
        claybot.Ore, _ = strconv.Atoi(caps[2])

        obsbot := BotCost{}
        obsbot.Ore, _ = strconv.Atoi(caps[3])
        obsbot.Clay, _ = strconv.Atoi(caps[4])

        geobot := BotCost{}
        geobot.Ore, _ = strconv.Atoi(caps[5])
        geobot.Obsidian, _ = strconv.Atoi(caps[6])

        bp := Blueprint{
            Id: id,
            OreBotCost: orebot,
            ClayBotCost: claybot,
            ObsidianBotCost: obsbot,
            GeodeBotCost: geobot,
            Inventory : BotCost{},

            nOreBots: 1,
            nClayBots: 0,
            nObsidianBots: 0,
            nGeodeBots: 0,
        }

        blueprints = append(blueprints, &bp)
    }

    return blueprints, nil
}
