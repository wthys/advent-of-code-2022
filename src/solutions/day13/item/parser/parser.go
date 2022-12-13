package parser

import (
    "io"
    "fmt"
    "bufio"
    "bytes"
    "strconv"
    "strings"

    "github.com/wthys/advent-of-code-2022/solutions/day13/item"
)

type (
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
    s.r.UnreadRune()
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

func (p *Parser) Parse() (*item.Item, error) {
    item, err := p.parseItem(0)
    if err != nil {
        return nil, err
    }

    return item, nil
}

func (p *Parser) parseItem(depth int) (*item.Item, error) {
    tok, lit := p.scanIgnoreWhiteSpace()

    if false {
        sep := strings.Repeat("--", depth)
        fmt.Println(sep)
    }
    switch tok {
    case INTEGER:
        value, _ := strconv.Atoi(lit)
        return item.ItemValue(value), nil
    case LBRACKET:
        items := []item.Item{}
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
        return item.ItemList(items...), nil
    case EOF:
        return nil, nil
    default:
        return nil, fmt.Errorf(`found %q, expected integer or "["`, lit)
    }
}
