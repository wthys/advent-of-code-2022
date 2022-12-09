package util

import (
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
