package lexer

import (
	"log"
	"strings"
	"unicode"

	"github.com/Mixturka/AutomataForge/internal/lexer/token"
)

type Lexer struct {
	mappings map[rune]token.Token
}

func NewLexer() Lexer {
	return Lexer{
		mappings: map[rune]token.Token{
			'+': token.NewToken(token.Plus, '+', 3),
			'*': token.NewToken(token.Star, '*', 3),
			'|': token.NewToken(token.Alter, '|', 1),
			'.': token.NewToken(token.Concat, '.', 2),
			'?': token.NewToken(token.Optional, '?', 3),
			'(': token.NewToken(token.LParen, '(', 4),
			')': token.NewToken(token.RParen, ')', 4),
			'[': token.NewToken(token.LSquare, '[', 4),
			']': token.NewToken(token.RSquare, ']', 4),
			'-': token.NewToken(token.Dash, '-', 4),
		},
	}
}

// Splits provided string into slice of tokens
func (l Lexer) Tokenize(s string) []token.Token {
	var tokens []token.Token
	runes := []rune(s)
	for i := 0; i < len(s); i++ {
		r := runes[i]
		// if next rune is character/digit/grouping/class/etc explicitly show concatenation '.'
		if i > 0 && (unicode.IsLetter(r) || unicode.IsDigit(r) || r == '[' || r == '(') && (unicode.IsLetter(runes[i-1]) ||
			unicode.IsDigit(runes[i-1]) || runes[i-1] == ')' || runes[i-1] == '*' || runes[i-1] == '+' || runes[i-1] == '?') {
			tokens = append(tokens, l.mappings['.'])
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			tokens = append(tokens, token.NewToken(token.Char, r, 0))
			continue
		}

		tok, ok := l.mappings[r]
		if !ok {
			log.Fatalf("Invalid character in input: %c", r)
		}

		if tok.Type == token.LSquare {
			endIdx := strings.IndexRune(s[i:], ']') + i
			if endIdx == -1 {
				log.Fatalf("Unmatched opening bracket at position %d", i)
			}
			if i+1 == endIdx {
				i++
				continue
			}
			tokens = append(tokens, l.evaluateSymbolClass(s[i+1:endIdx])...)
			i = endIdx
			continue
		}
		tokens = append(tokens, tok)
	}

	return tokens
}

func (l Lexer) evaluateSymbolClass(s string) []token.Token {
	var tokens []token.Token
	runes := []rune(s)

	// if length is > 1 give highest priority
	if len(s) > 1 {
		tokens = append(tokens, l.mappings['('])
	}
	for i := 0; i < len(runes); i++ {
		log.Printf("Cur: %c\n", runes[i])
		if i < len(runes)-1 && runes[i+1] != '-' {
			tokens = append(tokens, token.NewToken(token.Char, runes[i], 0))
			tokens = append(tokens, l.mappings['|'])
			// if we have "\\-" increase i to add '-' character
			if runes[i+1] == '\\' {
				i++
			}
		} else if i < len(runes)-2 && runes[i+1] == '-' {
			startRune, endRune := runes[i], runes[i+2]

			if startRune > endRune {
				log.Fatal("start cannot be greater then end")
			}
			for ; startRune <= endRune; startRune++ {
				tokens = append(tokens, token.NewToken(token.Char, startRune, 0))
				tokens = append(tokens, l.mappings['|'])
			}
			i += 2
		} else {
			tokens = append(tokens, token.NewToken(token.Char, runes[i], 0))
			tokens = append(tokens, l.mappings['|'])
		}
	}
	tokens = tokens[:len(tokens)-1] // remove unnecessary '|' token in the end
	// Close opened paren
	if len(s) > 1 {
		tokens = append(tokens, l.mappings[')'])
	}

	return tokens
}
