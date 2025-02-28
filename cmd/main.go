package main

import (
	"fmt"

	"github.com/Mixturka/AutomataForge/internal/lexer"
	"github.com/Mixturka/AutomataForge/internal/nfa"
	"github.com/Mixturka/AutomataForge/internal/parser"
)

func main() {
	lexer := lexer.NewLexer()
	// regex := "[a-c]*c+"
	regex := "fee|fie"
	tokens := lexer.Tokenize(regex)
	rpn := parser.ParseToRPN(tokens)
	for _, token := range rpn {
		fmt.Println(token)
	}

	nfa := nfa.BuildNFA(rpn)
	fmt.Println(nfa.PrettyPrint())

	dfa := nfa.BuildDFA()
	fmt.Println(len(dfa.Accepts))
	dfa.PrettyPrint()
	dfa.Minimize()

	dfa.PrettyPrint()
}
