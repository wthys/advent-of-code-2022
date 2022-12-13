package item

import (
    "fmt"
    "strings"
)

type (
    Item struct {
        Value int
        List []Item
        IsList bool
    }
)

func ItemValue(value int) *Item {
    return &Item{Value: value, List: nil, IsList: false}
}

func ItemList(items ...Item) *Item {
    return &Item{Value: 0, List: items, IsList: true}
}

func (i Item) Cmp(o Item) int {
    if i.IsList && o.IsList {
        for idx, item := range i.List {
            if idx >= len(o.List) {
                return 1
            }
            res := item.Cmp(o.List[idx])
            if res != 0 {
                return res
            }
        }

        if len(o.List) > len(i.List) {
            return -1
        }

        return 0
    }

    if !i.IsList && !o.IsList {
        switch {
        case i.Value < o.Value:
            return -1
        case i.Value > o.Value:
            return 1
        default:
            return 0
        }
    }

    if !i.IsList {
        return ItemList(i).Cmp(o)
    }

    return i.Cmp(*ItemList(o))
}

func (i Item) String() string {
    if i.IsList {
        str := strings.Builder{}
        fmt.Fprint(&str, "[")
        first := true
        for _, item := range i.List {
            if !first {
                fmt.Fprint(&str, ",")
            }
            fmt.Fprint(&str, item.String())
            first = false
        }
        fmt.Fprint(&str, "]")
        return str.String()
    }
    return fmt.Sprint(i.Value)
}
