package day7

import (
    "fmt"
    "strings"
)

const (
    Parent = ".."
)

type (
    Command struct{
        Name string
        Argument string
    }
    Node interface{
        Name() string
        Size() int
    }
    Directory struct {
        name string
        parent *Directory
        nodes []Node
    }
    File struct {
        name string
        size int
    }
)


func (d Directory) Name() string {
    return d.name
}

func (d Directory) Size() int {
    total := 0
    for _, node := range d.nodes {
        if node != nil {
            total += node.Size()
        }
    }
    return total
}

func (d Directory) GetNode(name string) (Node, error) {
    for _, node := range d.nodes {
        if node.Name() == name {
            return node, nil
        }
    }
    return nil, fmt.Errorf("could not find %q", name)
}

func (d *Directory) AddDir(name string) (*Directory, error) {

    node, err := d.GetNode(name)
    if err == nil {
        dir, ok := node.(*Directory)
        if ok {
            return dir, nil
        }
        return nil, fmt.Errorf("cannot add dir: %q already taken", name)
    }

    dir := &Directory{name, d, []Node{}}
    d.nodes = append(d.nodes, dir)

    return dir, nil
}

func (d *Directory) AddFile(name string, size int) (*File, error) {
    node, err := d.GetNode(name)
    if err == nil {
        file, ok := node.(File)
        if ok {
            return &file, nil
        }
        return nil, fmt.Errorf("cannot add file: %q already taken", name)
    }

    file := &File{name, size}
    d.nodes = append(d.nodes, file)

    return file, nil
}

func (d Directory) List() []Node {
    return d.nodes
}

func (d Directory) Parent() *Directory {
    return d.parent
}

func (d Directory) Path() string {
    p := d.Parent()
    if p != nil {
        return fmt.Sprintf("%v/%v", p.Path(), d.Name())
    }

    return ""
}

func (d Directory) String() string {
    return fmt.Sprintf("dir %v", d.Name())
}



func (f File) Name() string {
    return f.name
}

func (f File) Size() int {
    return f.size
}

func (f File) String() string {
    return fmt.Sprintf("%v %v", f.size, f.name)
}



func TreePrint(dir *Directory) {

    printNode := func(node Node, depth int) {
        dir := GetDir(node)
        indent := strings.Repeat(" ", depth)
        if dir != nil {
            fmt.Printf("%s- %s (dir, size=%d)\n", indent, dir.Name(), dir.Size())
            return
        }
        fmt.Printf("%s- %s (file, size=%d)\n", indent, node.Name(), node.Size())

    }

    TreeWalk(dir, 0, printNode)

}


func GetDir(node Node) *Directory {
    switch v := node.(type) {
    case *Directory:
        return v
    case Directory:
        return &v
    default:
        return nil
    }
}


func TreeWalk(node Node, depth int, action func(Node, int)) {
    action(node, depth)

    if dir := GetDir(node); dir != nil {
        for _, node := range dir.List() {
            TreeWalk(node, depth+1, action)
        }
    }
}
