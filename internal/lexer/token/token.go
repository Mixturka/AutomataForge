package token

import (
	"fmt"
)

type TokenType int

const (
	Char     TokenType = iota
	RParen             // )
	LParen             // (
	RSquare            // ]
	LSquare            // [
	Dash               // -
	Star               // *
	Plus               // +
	Optional           // ?
	Alter              // |
	Concat             // .
	Eps                // epsilon transition
)

type Token struct {
	Type       TokenType
	Value      rune
	Precedence int
}

func NewToken(t TokenType, v rune, precedence int) Token {
	return Token{
		Type:       t,
		Value:      v,
		Precedence: precedence,
	}
}

func IsOperator(t Token) bool {
	switch t.Type {
	case Char, LParen, RParen:
		return false
	default:
		return true
	}
}

func (t TokenType) String() string {
	switch t {
	case Char:
		return "Char"
	case RParen:
		return "RParen"
	case LParen:
		return "LParen"
	case RSquare:
		return "RSquare"
	case LSquare:
		return "LSquare"
	case Dash:
		return "Dash"
	case Star:
		return "Star"
	case Plus:
		return "Plus"
	case Optional:
		return "Optional"
	case Alter:
		return "Alter"
	case Concat:
		return "Concat"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Value: '%c'}", t.Type, t.Value)
}
