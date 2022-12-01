package main

import (
    "errors"
    "fmt"
    "io"
)


const (
    unknown = "unknown"
    unsolved = "unsolved"
    undefined = "undefined"
    inProgress = "in progress"
)

var (
    NotImplemented = errors.New("Not implemented")
)

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

func Solve(solver Solver, input io.Reader) (Result, error) {
    res := Result{
        Name: solver.Day(),
        Part1: unsolved,
        Part2: unsolved,
    }

    lines, err := ReadLines(input)

    if err != nil {
        return Result{}, fmt.Errorf("failed to read: %w", err)
    }

    if err := res.AddAnswers(solver, lines); err != nil {
        return Result{}, fmt.Errorf("failed to add answers: %w", err)
    }

    return res, nil
}

func (r *Result) AddAnswers(s Solver, input []string) error {
    part1, err := s.Part1(input)
    if err != nil && !errors.Is(err, NotImplemented) {
        return fmt.Errorf("failed to solve Part1: %w", err)
    }

    part2, err := s.Part2(input)
    if err != nil && !errors.Is(err, NotImplemented) {
        return fmt.Errorf("failed to solve Part2: %w", err)
    }

    r.Part1 = part1
    r.Part2 = part2

    return nil
}
