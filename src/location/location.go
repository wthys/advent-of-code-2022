package location

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/wthys/advent-of-code-2022/util"
)

var (
    reFromStr = regexp.MustCompile("[(]\\s*(\\d+)\\s*,\\s*(\\d+)\\s*[)]")
)

func New(x, y int) Location {
    return Location{x, y}
}

func FromString(input string) (Location, error) {
    none := Location{}
    caps := reFromStr.FindStringSubmatch(input)
    if caps == nil {
        return none, fmt.Errorf("no Location representation found")
    }

    var (
        err error = nil
        x, y int
    )

    x, err = strconv.Atoi(caps[1])
    if err != nil {
        return none, err
    }
    y, err = strconv.Atoi(caps[2])
    if err != nil {
        return none, err
    }

    return New(x, y), nil
}


type Location struct {
    X, Y int
}

func (l Location) String() string {
    return fmt.Sprintf("(%d,%d)", l.X, l.Y)
}

func (l Location) Add(o Location) Location {
    return New(l.X + o.X, l.Y + o.Y)
}

func (l Location) Scale(scale int) Location {
    return New(l.X * scale, l.Y * scale)
}

func (l Location) Subtract(o Location) Location {
    return New(l.X - o.X, l.Y - o.Y)
}

func (l Location) Unit() Location {
    return New(util.Sign(l.X), util.Sign(l.Y))
}
