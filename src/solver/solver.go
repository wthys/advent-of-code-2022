package solver

import (
    "errors"
    "fmt"
    "io"
    "context"
    "time"
)


const (
    Unknown = "unknown"
    Unsolved = "unsolved"
    Undefined = "undefined"
    InProgress = "in progress"
)

var (
    ErrNotImplemented = errors.New("Not implemented")
)

func NotImplemented() (string, error) {
    return Unsolved, ErrNotImplemented
}

func Solved[T any](value T) (string, error) {
    return fmt.Sprintf("%v", value), nil
}

func Error(err error) (string, error) {
    return Unsolved, err
}


type Day int


type Solver interface{
    Part1(input []string) (string, error)
    Part2(input []string) (string, error)
    Day() string
}

var (
    solvers = make(map[string]Solver)
)


func Register(solver Solver) {
    if solver == nil {
        panic("puzzle: Register solver is nil")
    }

    name := solver.Day()

    if _, dup := solvers[name]; dup {
        panic(fmt.Errorf("puzzle: Register called twice for solver [%s]", name))
    }

    solvers[name] = solver
}

func GetSolver(day string) (Solver, error) {
    if day == "" {
        return nil, errors.New("empty puzzle day")
    }

    solver, exist := solvers[day]
    if !exist {
        return nil, fmt.Errorf("%s: %w", day, errors.New("unknown puzzle day"))
    }

    return solver, nil
}

func Solve(solver Solver, input io.Reader, ctx context.Context) (Result, error) {
    res := Result{
        Name: solver.Day(),
        Part1: Unsolved,
        Part2: Unsolved,
        Elapsed: nil,
    }

    lines, err := ReadLines(input)

    if err != nil {
        return Result{}, fmt.Errorf("failed to read: %w", err)
    }

    if err := res.AddAnswers(solver, lines, ctx); err != nil {
        return Result{}, fmt.Errorf("failed to add answers: %w", err)
    }

    return res, nil
}

func (r *Result) AddAnswers(s Solver, input []string, ctx context.Context) error {
    elapsed, ok := ctx.Value("elapsed").(bool)
    if !ok {
        elapsed = false
    }

    durations := []time.Duration{}

    var start time.Time

    if (elapsed) {
        start = time.Now()
    }
    part1, err := s.Part1(input)
    if (elapsed) {
        durations = append(durations, time.Since(start))
    }
    if err != nil && !errors.Is(err, ErrNotImplemented) {
        return fmt.Errorf("failed to solve Part1: %w", err)
    }

    if (elapsed) {
        start = time.Now()
    }
    part2, err := s.Part2(input)
    if (elapsed) {
        durations = append(durations, time.Since(start))
    }
    if err != nil && !errors.Is(err, ErrNotImplemented) {
        return fmt.Errorf("failed to solve Part2: %w", err)
    }

    if !elapsed {
        durations = nil
    }

    r.Part1 = part1
    r.Part2 = part2
    r.Elapsed = durations

    return nil
}
