package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Mixturka/AutomataForge/internal/config"
	"github.com/Mixturka/AutomataForge/internal/dfa"
	"github.com/Mixturka/AutomataForge/internal/nfa"
)

type CliConfig struct {
	InputCfg   string
	outputPath string
}

type TableOutput struct {
	ClassifierTable map[rune]int          `json:"classifierTable"`
	TransitionTable [][]int               `json:"transitionTable"`
	TokenTypeTable  map[int]dfa.TokenInfo `json:"tokenTypeTable"`
}

func parseFlags() *CliConfig {
	var cfg CliConfig
	flag.StringVar(&cfg.outputPath, "o", "", "Output file (don't write this flag for stdout)")
	flag.StringVar(&cfg.outputPath, "output", "-", "Output file (alias)")
	flag.StringVar(&cfg.InputCfg, "i", "config.yml", `"Config path. Default - "config.yml`)
	flag.StringVar(&cfg.InputCfg, "input", "config.yml", `"Config path. Default - "config.yml`)
	flag.Parse()

	return &cfg
}

func main() {
	stateIdGen := nfa.NewStateIdGenerator()
	unifiedNfa := nfa.NewUnifiedNfa(&stateIdGen)

	cliFlags := parseFlags()
	tokenConfigs, err := config.ParseConfig(cliFlags.InputCfg)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err.Error())
	}

	for _, tokenCfg := range tokenConfigs {
		tokenNfa := nfa.BuildNFA(&stateIdGen, tokenCfg.Pattern, tokenCfg.Name, tokenCfg.Priority)
		unifiedNfa.AddRegex(tokenNfa)
	}

	dfa := unifiedNfa.BuildDFA()
	dfa.Minimize()

	classifierTable := dfa.BuildClassifierTable()
	transitionTable := dfa.BuildTransitionTable(classifierTable)
	tokenTypeTable := dfa.BuildTypeTable()
	tables := TableOutput{
		ClassifierTable: classifierTable,
		TransitionTable: transitionTable,
		TokenTypeTable:  tokenTypeTable,
	}

	encodedTables := make([]byte, 0)
	if encodedTables, err = json.Marshal(tables); err != nil {
		log.Fatalf("Failed to encode tables into JSON format: %v", err)
	}
	if err := os.WriteFile(cliFlags.outputPath, encodedTables, 0644); err != nil {
		log.Fatalf("Failed to write json into file %s: %v", cliFlags.outputPath, err)
	}
}
