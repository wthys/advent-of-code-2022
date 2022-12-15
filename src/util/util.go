package util

import (
    "fmt"
    "golang.org/x/exp/constraints"
)

type Number interface {
    constraints.Integer | constraints.Float
}

func Sign[T Number](val T) T {
    if val == 0 {
        return 0
    }
    return val / Abs(val)
}

func Abs[T Number](val T) T {
    if val < T(0) {
        return -val
    }
    return val
}

func Humanize(val int) string {
    if val < 1000 {
        return fmt.Sprint(val)
    }

    val = val / 1000
    if val < 1000 {
        return fmt.Sprintf("%vK", val)
    }

    val = val / 1000
    if val < 1000 {
        return fmt.Sprintf("%vM", val)
    }

    return fmt.Sprintf("%vG", val/1000)
}
