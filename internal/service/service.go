package service

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"unicode"
)

const Precision = 256

// Token types
const (
	NUMBER = "NUMBER"
	PLUS   = "PLUS"   // +
	MINUS  = "MINUS"  // -
	TIMES  = "TIMES"  // *
	DIVIDE = "DIVIDE" // /
	POWER  = "POWER"  // ^
	FACTOR = "FACTOR" // !
	LPAREN = "LPAREN" // (
	RPAREN = "RPAREN" // )
	EOF    = "EOF"
)

type Token struct {
	Type  string
	Value string
}

type Lexer struct {
	input  string
	pos    int
	tokens []Token
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.tokenize()
	return l
}

func (l *Lexer) tokenize() error {
	var tokens []Token
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case unicode.IsDigit(rune(ch)) || ch == '.':
			tokens = append(tokens, Token{Type: NUMBER, Value: l.readNumber()})
		case ch == '+':
			tokens = append(tokens, Token{Type: PLUS, Value: "+"})
			l.pos++
		case ch == '-':
			tokens = append(tokens, Token{Type: MINUS, Value: "-"})
			l.pos++
		case ch == '*':
			tokens = append(tokens, Token{Type: TIMES, Value: "*"})
			l.pos++
		case ch == '/':
			tokens = append(tokens, Token{Type: DIVIDE, Value: "/"})
			l.pos++
		case ch == '^':
			tokens = append(tokens, Token{Type: POWER, Value: "^"})
			l.pos++
		case ch == '(':
			tokens = append(tokens, Token{Type: LPAREN, Value: "("})
			l.pos++
		case ch == ')':
			tokens = append(tokens, Token{Type: RPAREN, Value: ")"})
			l.pos++
		case ch == '!':
			tokens = append(tokens, Token{Type: FACTOR, Value: "!"})
			l.pos++
		case ch == ' ' || ch == '\t':
			l.pos++
		default:
			return fmt.Errorf("invalid character: '%c'", ch)
		}
	}
	tokens = append(tokens, Token{Type: EOF, Value: ""})
	l.tokens = tokens

	return nil
}

func (l *Lexer) readNumber() string {
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(l.input[l.pos:])
	l.pos += len(match)
	return match
}

type Parser struct {
	lexer   *Lexer
	current Token
}

func NewParser(input string) *Parser {
	l := NewLexer(input)
	p := &Parser{lexer: l}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	if len(p.lexer.tokens) > 0 {
		p.current = p.lexer.tokens[0]
		p.lexer.tokens = p.lexer.tokens[1:]
	}
}

func (p *Parser) parseExpression() (*big.Float, error) {
	result, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.current.Type == PLUS || p.current.Type == MINUS {
		if p.current.Type == PLUS {
			p.nextToken()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			result = new(big.Float).Add(result, right)
		} else if p.current.Type == MINUS {
			p.nextToken()
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			result = new(big.Float).Sub(result, right)
		}
	}
	return result, nil
}

func (p *Parser) parseTerm() (*big.Float, error) {
	result, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.current.Type == TIMES || p.current.Type == DIVIDE {
		if p.current.Type == TIMES {
			p.nextToken()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			result = new(big.Float).Mul(result, right)
		} else if p.current.Type == DIVIDE {
			p.nextToken()
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			if right.Cmp(big.NewFloat(0)) == 0 {
				panic("integer division by zero")
			}
			result = new(big.Float).Quo(result, right)
		}
	}
	return result, nil
}

func (p *Parser) parseFactor() (*big.Float, error) {
	result, err := p.parsePower()
	if err != nil {
		return nil, err
	}
	for p.current.Type == FACTOR {
		p.nextToken()
		result, err = factorialFloat(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (p *Parser) parsePower() (*big.Float, error) {
	result, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	if p.current.Type == POWER {
		p.nextToken()
		exponent, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		result, err = floatPow(result, exponent)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (p *Parser) parsePrimary() (*big.Float, error) {
	if p.current.Type == NUMBER {
		f, ok := big.NewFloat(0).SetPrec(Precision).SetString(p.current.Value)
		if !ok {
			return nil, fmt.Errorf("error parsing number: %v", p.current.Value)
		}
		p.nextToken()
		return f, nil
	} else if p.current.Type == MINUS {
		p.nextToken()
		result, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		return new(big.Float).Neg(result), nil
	} else if p.current.Type == LPAREN {
		p.nextToken()
		result, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.current.Type != RPAREN {
			return nil, fmt.Errorf("expecting closing parenthesis")
		}
		p.nextToken()
		return result, nil
	} else {
		return nil, fmt.Errorf("unexpected token: %s", p.current.Type)
	}
}

func floatPow(base, exp *big.Float) (*big.Float, error) {
	if isInteger(exp) {
		intExp, _ := exp.Int64()
		return powIntExp(base, intExp), nil
	}

	return nil, fmt.Errorf("unexpected expression")
}

func powIntExp(base *big.Float, exp int64) *big.Float {
	result := big.NewFloat(1).SetPrec(Precision)
	baseCopy := new(big.Float).Set(base)

	if exp == 0 {
		return result
	}
	if exp < 0 {
		baseCopy = new(big.Float).Quo(big.NewFloat(1).SetPrec(Precision), baseCopy)
		exp = -exp
	}

	for exp > 0 {
		if exp&1 == 1 {
			result = new(big.Float).Mul(result, baseCopy)
		}
		baseCopy = new(big.Float).Mul(baseCopy, baseCopy)
		exp >>= 1
	}
	return result
}

func factorialFloat(f *big.Float) (*big.Float, error) {
	if !isInteger(f) {
		return nil, fmt.Errorf("factorial of non-integer number is not defined")
	}
	n, _ := f.Int64()
	if n < 0 {
		return nil, fmt.Errorf("factorial of negative number is not defined")
	}

	result := big.NewFloat(1).SetPrec(Precision)
	for i := int64(1); i <= n; i++ {
		result = new(big.Float).Mul(result, big.NewFloat(float64(i)))
	}
	return result, nil
}

func isInteger(f *big.Float) bool {
	_, acc := f.Int(nil)
	return acc == big.Exact
}

type CalcService struct{}

func NewCalcService() *CalcService {
	return &CalcService{}
}
func (c *CalcService) Evaluate(expr string) (*big.Float, error) {
	parser := NewParser(strings.TrimSpace(expr))
	result, err := parser.parseExpression()
	if err != nil {
		return nil, err

	}

	if parser.current.Type != EOF {
		return nil, fmt.Errorf("expecting end of expression, got: %s", parser.current.Type)
	}

	return result, nil
}
