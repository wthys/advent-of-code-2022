package blockoutput

import (
    "fmt"
    "strings"
    "crypto/md5"
)


type blockletter [][]rune

var str2letter = map[string]blockletter{}
var letter2str = map[[16]byte]string{}

const (
    on = '#'
    off = '.'
)


func FromBlockLetters(input [][]rune) string {
    letters := splitBlockLetters(input)

    str := strings.Builder{}

    for _, letter := range letters {
        val, ok := letter2str[letter.Hash()]
        if ok {
            fmt.Fprint(&str, val)
        } else {
            fmt.Fprint(&str, "?")
        }
    }

    return str.String()
}

func ToBlockLetters(value string) [][]rune {
    letters := []blockletter{}

    for _, char := range value {
        letter, ok := str2letter[string(char)]
        if ok {
            letters = append(letters, letter)
        }
    }

    return joinBlockLetters(letters)

}

func (bl blockletter) Hash() [16]byte {
    data := []byte{}
    for _, row := range bl {
        data = append(data, []byte(string(row))...)
    }
    return md5.Sum(data)
}

func splitBlockLetters(input [][]rune) []blockletter {
    letters := []blockletter{}

    total := len(input[0])/5

    for i := 0; i < total; i++ {
        letter := blockletter{}
        for _, row := range input {
            letter = append(letter, row[5*i:5*i+4])
        }

        letters = append(letters, letter)
    }

    return letters
}

func joinBlockLetters(letters []blockletter) [][]rune {
    ll := len(letters)*5

    output := [][]rune{}
    for j := 0; j < 6; j++ {
        output = append(output, []rune(strings.Repeat(string(off), ll)))
    }

    for i, letter := range letters {
        for j, _ := range output {
            for k := range []int{0,1,2,3} {
                output[j][5*i+k] = letter[j][k]
            }
        }
    }

    return output
}

func init() {

    //TODO verify A
    str2letter["A"] = [][]rune{
        []rune(".##."),
        []rune("#..#"),
        []rune("####"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
    }

    str2letter["B"] = [][]rune{
        []rune("###."),
        []rune("#..#"),
        []rune("###."),
        []rune("#..#"),
        []rune("#..#"),
        []rune("###."),
    }

    str2letter["C"] = [][]rune{
        []rune(".##."),
        []rune("#..#"),
        []rune("#..."),
        []rune("#..."),
        []rune("#..#"),
        []rune(".##."),
    }

    //TODO verify D
    str2letter["D"] = [][]rune{
        []rune("###."),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("###."),
    }

    str2letter["E"] = [][]rune{
        []rune("####"),
        []rune("#..."),
        []rune("###."),
        []rune("#..."),
        []rune("#..."),
        []rune("####"),
    }

    str2letter["F"] = [][]rune{
        []rune("####"),
        []rune("#..."),
        []rune("###."),
        []rune("#..."),
        []rune("#..."),
        []rune("#..."),
    }

    str2letter["G"] = [][]rune{
        []rune(".##."),
        []rune("#..#"),
        []rune("#..."),
        []rune("#.##"),
        []rune("#..#"),
        []rune(".###"),
    }

    str2letter["H"] = [][]rune{
        []rune("#..#"),
        []rune("#..#"),
        []rune("####"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
    }

    // I

    str2letter["J"] = [][]rune{
        []rune("..##"),
        []rune("...#"),
        []rune("...#"),
        []rune("...#"),
        []rune("#..#"),
        []rune(".##."),
    }

    str2letter["K"] = [][]rune{
        []rune("#..#"),
        []rune("#.#."),
        []rune("##.."),
        []rune("#.#."),
        []rune("#.#."),
        []rune("#..#"),
    }

    str2letter["L"] = [][]rune{
        []rune("#..."),
        []rune("#..."),
        []rune("#..."),
        []rune("#..."),
        []rune("#..."),
        []rune("####"),
    }

    // M

    // N

    //TODO verify O
    str2letter["O"] = [][]rune{
        []rune(".##."),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune(".##."),
    }

    str2letter["P"] = [][]rune{
        []rune("###."),
        []rune("#..#"),
        []rune("#..#"),
        []rune("###."),
        []rune("#..."),
        []rune("#..."),
    }

    // Q

    str2letter["R"] = [][]rune{
        []rune("###."),
        []rune("#..#"),
        []rune("#..#"),
        []rune("###."),
        []rune("#.#."),
        []rune("#..#"),
    }

    //TODO verify S
    str2letter["S"] = [][]rune{
        []rune(".##."),
        []rune("#..#"),
        []rune(".#.."),
        []rune("..#."),
        []rune("#..#"),
        []rune(".##."),
    }

    // T

    str2letter["U"] = [][]rune{
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune("#..#"),
        []rune(".##."),
    }

    // V

    // W

    // X

    // Y

    str2letter["Z"] = [][]rune{
        []rune("####"),
        []rune("...#"),
        []rune("..#."),
        []rune(".#.."),
        []rune("#..."),
        []rune("####"),
    }

    // Reverse the map
    for key, letter := range str2letter {
        letter2str[letter.Hash()] = key
    }
}
