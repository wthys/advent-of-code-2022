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

func IIf[T any](condition bool, yes, no T) T {
    if condition {
        return yes
    }
    return no
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

func PermutationDo[T any](k int, values []T, doer func (permutation []T)) {
    c := []int{}
    for i := 0; i < k; i++ {
        c = append(c, 0)
    }

    array := values

    doer(array)

    i := 1
    for i < k {
        if c[i] >= i {
            c[i] = 0
            i += 1
            continue
        }

        if i % 2 == 0 {
            array[0], array[i] = array[i], array[0]
        } else {
            array[c[i]], array[i] = array[i], array[c[i]]
        }
        doer(array)

        c[i] += 1
        i = 1
    }
}

func PairWiseDo[T any](values []T, doer func(a, b T)) {
    if len(values) < 2 {
        return
    }

    prev := values[0]
    for _, val := range values[1:] {
        doer(prev, val)
        prev = val
    }
}

func Max[T Number](values ...T) T {
    if len(values) == 0 {
        panic("need at least one value")
    }
    best := values[0]

    if len(values) == 1 {
        return best
    }

    for _, value := range values[1:] {
        if value > best {
            best = value
        }
    }

    return best
}

func Min[T Number](values ...T) T {
    if len(values) == 0 {
        panic("need at least one value")
    }
    best := values[0]

    if len(values) == 1 {
        return best
    }

    for _, value := range values[1:] {
        if value < best {
            best = value
        }
    }

    return best
}
