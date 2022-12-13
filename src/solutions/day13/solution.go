package day13

import (
    "strings"
    "io"
    "bufio"
    "strconv"
    "bytes"
    "fmt"
    "sort"

    "github.com/wthys/advent-of-code-2022/solver"
)

type solution struct{}

func init() {
    solver.Register(solution{})
}

func (s solution) Day() string {
    return "13"
}

func (s solution) Part1(input []string) (string, error) {
    pairs, err := parseInput(input)
    if err != nil {
        return "", err
    }

    total := 0
    for i, pair := range pairs {
        if pair.inOrder() {
            total += i + 1
        }
    }

    return strconv.Itoa(total), nil
}

func (s solution) Part2(input []string) (string, error) {
    pairs, err := parseInput(input)
    if err != nil {
        return "", err
    }

    items := []Item{}

    for _, pair := range pairs {
        items = append(items, pair.Left, pair.Right)
    }

    sep2 := *ItemList(*ItemValue(2))
    sep6 := *ItemList(*ItemValue(6))

    items = append(items, sep2, sep6)

    sort.Sort(ItemSlice(items))

    decoderKey := 1
    for i, item := range items {
        if item.Cmp(sep2) == 0 || item.Cmp(sep6) == 0 {
            decoderKey *= i + 1
        }
    }

    return strconv.Itoa(decoderKey), nil
}

func (s ItemSlice) Swap(i, j int) {
    s[j], s[i] = s[i], s[j]
}

func (s ItemSlice) Len() int {
    return len(s)
}

func (s ItemSlice) Less(i, j int) bool {
    return s[i].Cmp(s[j]) < 0
}

func (p Pair) inOrder() bool {
    return p.Left.Cmp(p.Right) <= 0
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

func (i *Item) String() string {
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

func parseInput(input []string) ([]Pair, error) {

    pairs := []Pair{}

    var (
        left *Item = nil
    )

    for n, line := range input {
        if len(strings.TrimSpace(line)) == 0 {
            continue
        }

        item, err := NewParser(strings.NewReader(line)).Parse()
        if err != nil {
            return nil, fmt.Errorf("error on line %v: %v", n+1, err)
        }

        if left == nil {
            //fmt.Printf("left: %v\n", item)
            left = item
        } else {
            //fmt.Printf("right: %v\n", item)
            pairs = append(pairs, Pair{*left, *item})
            left = nil
        }
    }

    return pairs, nil
}

type (
    Pair struct {
        Left Item
        Right Item
    }

    Item struct {
        Value int
        List []Item
        IsList bool
    }

    ItemSlice []Item

    Token int

    Scanner struct {
        r *bufio.Reader
    }

    Parser struct {
        s *Scanner
        buf struct {
            tok Token
            lit string
            isUnscanned bool
        }
    }
)

const (
    ILLEGAL Token = iota
    WS
    INTEGER
    LBRACKET
    RBRACKET
    COMMA
    EOF
)

func ItemValue(value int) *Item {
    return &Item{Value: value, List: nil, IsList: false}
}

func ItemList(items ...Item) *Item {
    return &Item{Value: 0, List: items, IsList: true}
}

func NewScanner(r io.Reader) *Scanner {
    return &Scanner{r: bufio.NewReader(r)}
}

func NewParser(r io.Reader) *Parser {
    return &Parser{s: NewScanner(r)}
}

func isWhiteSpace(ch rune) bool {
    return ch == ' ' || ch == '\t'
}

func isDigit(ch rune) bool {
    return (ch >= '0' && ch <= '9')
}

func (s *Scanner) read() rune {
    ch, _, err := s.r.ReadRune()
    if err != nil {
        return rune(0)
    }
    return ch
}

func (s *Scanner) unread() {
    _ = s.r.UnreadRune()
}

func (s *Scanner) Scan() (tok Token, lit string) {
    ch := s.read()

    if isWhiteSpace(ch) {
        s.unread()
        return s.scanWhiteSpace()
    } else if isDigit(ch) {
        s.unread()
        return s.scanDigit()
    }

    switch ch {
    case '[':
        return LBRACKET, string(ch)
    case ']':
        return RBRACKET, string(ch)
    case ',':
        return COMMA, string(ch)
    case rune(0):
        return EOF, ""
    default:
        return ILLEGAL, string(ch)
    }
}

func (s *Scanner) scanWhiteSpace() (tok Token, lit string) {
    var buf bytes.Buffer
    buf.WriteRune(s.read())

    for {
        ch := s.read()
        if ch == rune(0) {
            break
        }

        if !isWhiteSpace(ch) {
            s.unread()
            break
        }

        buf.WriteRune(ch)
    }

    return WS, buf.String()
}

func (s *Scanner) scanDigit() (Token, string) {
    var buf bytes.Buffer
    buf.WriteRune(s.read())

    for {
        ch := s.read()
        if ch == rune(0) {
            break
        }

        if !isDigit(ch) {
            s.unread()
            break
        }

        buf.WriteRune(ch)
    }

    return INTEGER, buf.String()
}

func (p *Parser) scan() (Token, string) {
    if p.buf.isUnscanned {
        p.buf.isUnscanned = false
        return p.buf.tok, p.buf.lit
    }

    p.buf.tok, p.buf.lit = p.s.Scan()

    return p.buf.tok, p.buf.lit
}

func (p *Parser) unscan() {
    p.buf.isUnscanned = true
}

func (p *Parser) scanIgnoreWhiteSpace() (Token, string) {
    tok, lit := p.scan()
    if tok == WS {
        tok, lit = p.scan()
    }
    return tok, lit
}

func (p *Parser) Parse() (*Item, error) {
    item, err := p.parseItem(0)
    if err != nil {
        return nil, err
    }

    return item, nil
}

func (p *Parser) parseItem(depth int) (*Item, error) {
    tok, lit := p.scanIgnoreWhiteSpace()

    if false {
        sep := strings.Repeat("--", depth)
        fmt.Println(sep)
    }
    switch tok {
    case INTEGER:
        value, _ := strconv.Atoi(lit)
        return ItemValue(value), nil
    case LBRACKET:
        items := []Item{}
        for {
            tok, lit := p.scanIgnoreWhiteSpace()
            if tok == RBRACKET {
                break
            }
            if tok != COMMA {
                p.unscan()
            }
            item, err := p.parseItem(depth + 1)
            if err != nil {
                return nil, err
            }
            if item == nil {
                break
            }
            items = append(items, *item)

            tok, lit = p.scanIgnoreWhiteSpace()
            if tok == RBRACKET {
                break
            }
            switch tok {
            case COMMA:
                continue
            default:
                return nil, fmt.Errorf(`found %q, expected "]" or ","`, lit)
            }
        }
        return ItemList(items...), nil
    case EOF:
        return nil, nil
    default:
        return nil, fmt.Errorf(`found %q, expected integer or "["`, lit)
    }
}
