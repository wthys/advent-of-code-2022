package solver

import (
    "fmt"
    "io"
    "strings"

    "text/tabwriter"
)

type Result struct{
    Name string
    Part1 string
    Part2 string
}

func (r Result) String() string {
    if r.Part1 == "" {
        r.Part1 = Unsolved
    }

    if r.Part2 == "" {
        r.Part2 = Unsolved
    }

    if r.Name == "" {
        r.Name = Unknown
    }


    const linesnum = 10

    table := make([][]string, 0, linesnum)

    var empty []string

    header := []string{fmt.Sprintf("%s puzzle answer", r.Name)}
    part1 := []string{"part1", r.Part1}
    part2 := []string{"part2", r.Part2}

    table = append(table, empty, header, part1, part2, empty, empty)

    var buf strings.Builder
    if err := printTable(&buf, table); err != nil {
        panic(err)
    }

    content := strings.TrimSpace(buf.String())
    content = fmt.Sprintf("\n  %s\n", content)

    return content
}

func printTable(w io.Writer, table [][]string) error {
    const padding = 3

    writer := tabwriter.NewWriter(w, 0, 0, padding, ' ', tabwriter.DiscardEmptyColumns)

    var err error

    for _, line := range table {
        switch len(line) {
        case 0:
            _, err = fmt.Fprintf(w, "\n")
        case 1:
            _, err = fmt.Fprintf(writer, "\t" + strings.Join(line, "\t") + "\t\n")
        default:
            _, err = fmt.Fprintf(writer, "\t   " + strings.Join(line, "\t") + "\t\n")
        }

        if err != nil {
            return fmt.Errorf("fprintln: %w", err)
        }
    }

    if err := writer.Flush(); err != nil {
        return fmt.Errorf("flush: %w", err)
    }

    return nil
}
