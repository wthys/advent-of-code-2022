package location

import (
    "fmt"
    "regexp"
    "strconv"

    "github.com/wthys/advent-of-code-2022/util"
)

var (
    reFromStr = regexp.MustCompile("^\\s*[(]\\s*(-?\\d+)\\s*,\\s*(-?\\d+)\\s*[)]\\s*$")
    reFromStr3 = regexp.MustCompile("^\\s*[(]\\s*(-?\\d+)\\s*,\\s*(-?\\d+)\\s*,\\s*(-?\\d+)\\s*[)]\\s*$")
    ErrWrongFormat = fmt.Errorf("wrong Location format")
)

func New(x, y int) Location {
    return Location{x, y}
}

func New3(x, y, z int) Location3 {
    return Location3{x, y, z}
}

func FromString(input string) (Location, error) {
    none := Location{}
    caps := reFromStr.FindStringSubmatch(input)
    if caps == nil {
        return none, ErrWrongFormat
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

func FromString3(input string) (Location3, error) {
    none := Location3{}
    caps := reFromStr3.FindStringSubmatch(input)
    if caps == nil {
        return none, ErrWrongFormat
    }

    var (
        err error = nil
        x, y, z int
    )

    x, err = strconv.Atoi(caps[1])
    if err != nil {
        return none, err
    }

    y, err = strconv.Atoi(caps[2])
    if err != nil {
        return none, err
    }

    z, err = strconv.Atoi(caps[3])
    if err != nil {
        return none, err
    }

    return New3(x, y, z), nil
}


type (
    Location struct {
        X, Y int
    }

    Location3 struct {
        X, Y, Z int
    }
)

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

func (l Location) Manhattan() int {
    return util.Abs(l.X) + util.Abs(l.Y)
}

func (l Location3) String() string {
    return fmt.Sprintf("(%d,%d,%d)", l.X, l.Y, l.Z)
}

func (l Location3) Add(o Location3) Location3 {
    return New3(l.X + o.X, l.Y + o.Y, l.Z + o.Z)
}

func (l Location3) Scale(scale int) Location3 {
    return New3(l.X * scale, l.Y * scale, l.Z * scale)
}

func (l Location3) Subtract(o Location3) Location3 {
    return New3(l.X - o.X, l.Y - o.Y, l.Z - o.Z)
}

func (l Location3) Unit() Location3 {
    return New3(util.Sign(l.X), util.Sign(l.Y), util.Sign(l.Z))
}

func (l Location3) Manhattan() int {
    return util.Abs(l.X) + util.Abs(l.Y) + util.Abs(l.Z)
}
