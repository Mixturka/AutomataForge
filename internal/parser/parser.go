package parser

import (
	"github.com/Mixturka/AutomataForge/internal/lexer/token"
	"github.com/Mixturka/AutomataForge/pkg/stack"
)

// Simple shunting yard algorithm to reorder tokens according to Reverse Polish Notation (RPN)
func ParseToRPN(tokens []token.Token) []token.Token {
	stack := stack.NewStack[token.Token]()
	var rpn []token.Token

	for _, tok := range tokens {
		if tok.Type == token.Char {
			rpn = append(rpn, tok)
		} else if token.IsOperator(tok) {
			for stack.Length() > 0 && stack.Top().Precedence >= tok.Precedence && stack.Top().Type != token.LParen {
				rpn = append(rpn, stack.Pop())
			}
			stack.Push(tok)
		} else if tok.Type == token.LParen {
			stack.Push(tok)
		} else if tok.Type == token.RParen {
			for stack.Length() > 0 && stack.Top().Type != token.LParen {
				rpn = append(rpn, stack.Pop())
			}
			stack.Pop()
		}
	}

	for stack.Length() > 0 {
		rpn = append(rpn, stack.Pop())
	}
	return rpn
}
