package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Mixturka/AutomataForge/internal/config"
	"github.com/Mixturka/AutomataForge/internal/lexer"
	"github.com/Mixturka/AutomataForge/internal/parser"
)

func main() {
	configPath := flag.String("config", "automata-forge.yml", "--config=/path/to/config")
	flag.Parse()

	tokenConfig, err := config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err.Error())
	}

	fmt.Println(tokenConfig)

	lexer := lexer.NewLexer()
	// regex := "[a-c]*c+"
	regex := "fee|fie"
	tokens := lexer.Tokenize(regex)
	rpn := parser.ParseToRPN(tokens)
	for _, token := range rpn {
		fmt.Println(token)
	}

	// nfa := nfa.BuildNFA(rpn)
	// fmt.Println(nfa.PrettyPrint())

	// dfa := nfa.BuildDFA()
	// fmt.Println(len(dfa.Accepts))
	// dfa.PrettyPrint()
	// dfa.Minimize()

	// dfa.PrettyPrint()
}
