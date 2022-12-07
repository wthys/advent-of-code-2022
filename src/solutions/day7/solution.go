package day7

import (
    "fmt"
    "regexp"
    "strconv"
    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "7"
}

func parseInput(input []string) (*Directory, error) {

    root := &Directory{"/", nil, []Node{}}

    re_cmd := regexp.MustCompile("^[$]\\s+(cd|ls)(\\s+([a-zA-Z\\.]+|/))?\\s*$")
    re_dir := regexp.MustCompile("^dir\\s+([a-z\\.]+)\\s*$")
    re_file := regexp.MustCompile("^\\s*(\\d+)\\s+([a-z\\.]+)\\s*$")

    current := root

    for _, line := range input {
        cmd := re_cmd.FindStringSubmatch(line)
        if cmd != nil {
            switch cmd[1] {
            case "cd":
                switch cmd[3] {
                case Parent:
                    current = current.Parent()
                case "/":
                    current = root
                default:
                    dir, err := current.AddDir(cmd[3])
                    if err != nil {
                        return nil, err
                    }
                    current = dir
                }
            case "ls":
                continue
            default:
                return nil, fmt.Errorf("unknown command %q", cmd[0])
            }
        } else {
            dir := re_dir.FindStringSubmatch(line)
            if dir != nil {
                _, err := current.AddDir(dir[1])
                if err != nil {
                    return nil, err
                }
                continue
            }

            file := re_file.FindStringSubmatch(line)
            if file != nil {
                size, err := strconv.Atoi(file[1])
                if err != nil {
                    return nil, err
                }
                _, err = current.AddFile(file[2], size)
                if err != nil {
                    return nil, err
                }
                continue
            }
            return nil, fmt.Errorf("cannot use %q, help?", line)
        }
    }

    return root, nil

}


func (s solution) Part1(input []string) (string, error) {
    root, err := parseInput(input)
    if err != nil {
        return "", err
    }

    total := 0
    
    SumUnder100K := func(node Node, _ int) {
        if dir := GetDir(node); dir != nil {
            size := dir.Size()
            if size < 100_000 {
                //fmt.Printf("found dir %q: %v\n", dir.Name(), size)
                total += size
            }
        }
    }

    TreeWalk(root, 0, SumUnder100K)

    //TreePrint(root)

    return strconv.Itoa(total), nil
}

func (s solution) Part2(input []string) (string, error) {
    root, err := parseInput(input)
    if err != nil {
        return "", err
    }

    free := 70_000_000 - root.Size()
    needed := 30_000_000


    smallest := root.Size()
    //smallest_dir := root

    findDir := func(node Node, _ int) {
        if dir := GetDir(node); dir != nil {
            dirsize := dir.Size()
            if free + dirsize >= needed {
                if dirsize < smallest {
                    smallest = dirsize
                    //smallest_dir = dir
                }
            }
        }
    }

    TreeWalk(root, 0, findDir)

    //fmt.Printf("smallest dir is (%v, %v): %v\n", smallest_dir.Path(), smallest_dir, smallest)


    return strconv.Itoa(smallest), nil
}
